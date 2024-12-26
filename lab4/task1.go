package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"sync"
	"time"
)

// Ядро свёртки для размытия по Гауссу (3x3)
var gaussianKernel = [3][3]float64{
	{1 / 16.0, 1 / 8.0, 1 / 16.0},
	{1 / 8.0, 1 / 4.0, 1 / 8.0},
	{1 / 16.0, 1 / 8.0, 1 / 16.0},
}

func main() {
	// 1 задание
	handleChannels()

	// 2 задание
	// Выполняем обработку изображения обычным методом
	handleImageProcessing(false)

	// 3 задание
	// Выполняем обработку изображения с параллельной обработкой
	handleImageProcessing(true)

	fmt.Println("Программа завершена.")
	// 4 задание
	// Открываем изображение
	file, err := os.Open("input.png")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	// Декодируем изображение
	srcImg, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Ошибка при декодировании изображения:", err)
		return
	}

	// Преобразуем изображение в формат RGBA
	drawImg := toRGBA(srcImg)

	// Замеряем время выполнения
	start := time.Now()

	// Применяем фильтр с параллельной обработкой
	applyGaussianBlur(drawImg)

	duration := time.Since(start)
	fmt.Printf("Обработка изображения с матричным фильтром заняла: %v\n", duration)

	// Сохраняем результат
	outputFile, err := os.Create("output_blurred.png")
	if err != nil {
		fmt.Println("Ошибка при создании выходного файла:", err)
		return
	}
	defer outputFile.Close()

	// Сохраняем обработанное изображение
	err = png.Encode(outputFile, drawImg)
	if err != nil {
		fmt.Println("Ошибка при сохранении изображения:", err)
		return
	}

	fmt.Println("Изображение успешно обработано и сохранено в output_blurred.png")
}

// handleChannels запускает задачу работы с каналами
func handleChannels() {
	ch := make(chan int)
	var wg sync.WaitGroup

	// Увеличиваем счетчик WaitGroup
	wg.Add(1)

	// Запускаем функцию count в отдельной горутине
	go count(ch, &wg)

	// Отправляем числа в канал
	for i := 1; i <= 5; i++ {
		ch <- i
	}
	// Закрываем канал
	close(ch)

	// Ждем завершения горутины count
	wg.Wait()
}

// count читает числа из канала, возводит их в квадрат и выводит результат
func count(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for num := range ch {
		fmt.Printf("Квадрат числа %d: %d\n", num, num*num)
	}
}

// Преобразуем изображение в RGBA
// Преобразуем изображение в RGBA
func toRGBA(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	return rgba
}

// applyGaussianBlur применяет фильтр размытия по Гауссу к изображению
func applyGaussianBlur(img *image.RGBA) {
	bounds := img.Bounds()
	var wg sync.WaitGroup

	// Увеличиваем счетчик перед запуском горутин
	// Убедитесь, что горутины будут добавлены для всех строк (y-координаты)
	wg.Add(bounds.Max.Y - bounds.Min.Y - 2) // -2 для пропуска первого и последнего пикселя

	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		go func(y int) {
			defer wg.Done()

			for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
				// Применяем свёртку для каждого пикселя
				color := applyKernel(img, x, y)

				// Записываем результат в изображение
				img.Set(x, y, color)
			}
		}(y)
	}

	// Ожидаем завершения всех горутин
	wg.Wait()
}

func handleImageProcessing(parallel bool) {
	// Открываем файл изображения
	file, err := os.Open("input.png")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	// Декодируем изображение
	srcImg, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Ошибка при декодировании изображения:", err)
		return
	}

	// Преобразуем изображение в редактируемый формат
	drawImg := toRGBA(srcImg)

	// Замеряем время выполнения
	start := time.Now()

	if parallel {
		// Параллельная обработка
		filterParallel(drawImg)
	} else {
		// Обычная обработка
		filter(drawImg)
	}

	duration := time.Since(start)
	if parallel {
		fmt.Printf("Параллельная обработка заняла: %v\n", duration)
	} else {
		fmt.Printf("Обычная обработка заняла: %v\n", duration)
	}

	// Создаем новый файл для сохранения результата
	outputFile, err := os.Create(fmt.Sprintf("output_%v.png", parallel))
	if err != nil {
		fmt.Println("Ошибка при создании выходного файла:", err)
		return
	}
	defer outputFile.Close()

	// Сохраняем обработанное изображение
	err = png.Encode(outputFile, drawImg)
	if err != nil {
		fmt.Println("Ошибка при сохранении изображения:", err)
		return
	}

	fmt.Printf("Изображение успешно обработано и сохранено в output_%v.png\n", parallel)
}

// filter применяет преобразование к каждому пикселю изображения (обычный метод)
func filter(img draw.RGBA64Image) {
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Получаем цвет пикселя
			pixel := img.RGBA64At(x, y)

			// Вычисляем яркость (среднее значение каналов)
			gray := uint16((uint32(pixel.R) + uint32(pixel.G) + uint32(pixel.B)) / 3)

			// Создаем новый цвет в оттенках серого
			grayColor := color.RGBA64{R: gray, G: gray, B: gray, A: pixel.A}

			// Устанавливаем новый цвет пикселя
			img.SetRGBA64(x, y, grayColor)
		}
	}
}

// filterParallel применяет преобразование к одной строке пикселей изображения (параллельный метод)
func filterParallel(img draw.RGBA64Image) {
	bounds := img.Bounds()
	var wg sync.WaitGroup

	// Устанавливаем количество горутин равным количеству строк в изображении
	wg.Add(bounds.Max.Y - bounds.Min.Y)

	// Запускаем горутины для обработки каждой строки
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		go func(y int) {
			defer wg.Done() // Уменьшаем счетчик после завершения работы горутины

			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				// Получаем цвет пикселя
				pixel := img.RGBA64At(x, y)

				// Вычисляем яркость (среднее значение каналов)
				gray := uint16((uint32(pixel.R) + uint32(pixel.G) + uint32(pixel.B)) / 3)

				// Создаем новый цвет в оттенках серого
				grayColor := color.RGBA64{R: gray, G: gray, B: gray, A: pixel.A}

				// Устанавливаем новый цвет пикселя
				img.SetRGBA64(x, y, grayColor)
			}
		}(y)
	}

	// Ожидаем завершения всех горутин
	wg.Wait()
}

// applyKernel применяет ядро свёртки для одного пикселя
func applyKernel(img *image.RGBA, x, y int) color.Color {
	var r, g, b, a float64

	// Проходим по матрице 3x3 и применяем веса из ядра
	for ky := -1; ky <= 1; ky++ {
		for kx := -1; kx <= 1; kx++ {
			// Получаем цвет соседнего пикселя
			px := img.RGBAAt(x+kx, y+ky)

			// Умножаем значения на соответствующие веса из ядра
			r += float64(px.R) * gaussianKernel[ky+1][kx+1]
			g += float64(px.G) * gaussianKernel[ky+1][kx+1]
			b += float64(px.B) * gaussianKernel[ky+1][kx+1]
			a += float64(px.A) * gaussianKernel[ky+1][kx+1]
		}
	}

	// Возвращаем новый цвет после применения ядра
	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
}

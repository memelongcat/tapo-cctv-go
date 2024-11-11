package main //Название пакета

import ( //Импорт библиотек и модулей
	"encoding/json" //Работа с JSON
	"fmt"           // //Ввод-вывод
	"log"           //Логгирование ошибок
	"os"            //Работа с ФС
	"time"          //Время

	ffmpeg_go "github.com/u2takey/ffmpeg-go" //либа для работы с ffmpeg
)

func current_time() string { //получаем текущее время
	return time.Now().Format("15.04.05") //возвращаем время в формате HH.MM.SS
}

func current_date() string { //получаем текущую дату
	return time.Now().Format("02.01.2006") //возвращаем дату в DD.MM.YY
}

func file_directory(record_directory string, date string) string { //получаем название директpeории, создаем, если не существует (по моему, выглядит как говнокод, когда возвращает входное значение при удачном результате)
	directory := record_directory + date //директория для записей + текущая дата
	_, err := os.Stat(directory)         //кидаем в _ информацию о директории
	if os.IsNotExist(err) {              //если нет директории, создаем
		os.Mkdir(directory, 0755) //создаем директорию с уровнем доступа 0755
	}
	return directory
}

func get_ffmpeg_config(config_directory string) Config { //получаем JSON-файл с параметрами записи
	config_file, err := os.ReadFile(config_directory) //читаем JSON файл
	if err != nil {                                   //Логгируем ошибку при чтении файла и выходим из программы, если ошибка не nil
		log.Fatal(err)
	}
	var config Config                          //объявляем переменную config типа Config
	err = json.Unmarshal(config_file, &config) //анмаршалим JSON в структуру Config (TODO: прикутить обработчик ошибок)
	if err != nil {                            //Логгируем ошибку при анмаршалинге JSON и выходим из программы, если ошибка не nil
		log.Fatal(err)
	}
	return config
}

func cam_record(config Config, file_directory string, file_name string) {

	input := fmt.Sprintf("rtsp://%s:%s@%s/%s", config.Cam_user, config.Cam_password, config.Cam_ip, config.Cam_stream) //составляем параметр для input
	output := fmt.Sprintf("%s/%s.%s", file_directory, file_name, config.Filetype)                                      //составляем параметр для output

	err := ffmpeg_go.Input(input, ffmpeg_go.KwArgs{"t": config.Duration}).Output(output, ffmpeg_go.KwArgs{"vcodec": "copy", "b:v": config.Bitrate}).Run() //запускаем ffmpeg
	if err != nil {                                                                                                                                       //Логгируем ошибку при запуске ffmpeg и выходим из программы, если ошибка не nil
		log.Fatal(err)
	}

}

type Config struct { //структура JSON, аналогичная структуре файла
	Cam_ip           string `json:"cam_ip"`
	Cam_user         string `json:"cam_user"`
	Cam_password     string `json:"cam_password"`
	Cam_stream       string `json:"cam_stream"`
	Duration         string `json:"duration"`
	Filetype         string `json:"filetype"`
	Bitrate          string `json:"bitrate"`
	Record_directory string `json:"record_directory"`
}

func main() { //главная функция

	for {
		cfg := get_ffmpeg_config("./config.json")
		dir := file_directory(cfg.Record_directory, current_date())
		file := current_time()

		//fmt.Printf("%s", dir)
		cam_record(cfg, dir, file)
	}
}

// TODO: Сделать логгирование ошибок вместо отправки в _
package main //Название пакета

import ( //Импорт библиотек и модулей
	"encoding/json" //Работа с JSON
	"fmt"           // //Ввод-вывод
	"os"            //Работа с ФС
	"time"          //Время

	ffmpeg_go "github.com/u2takey/ffmpeg-go" //либа для работы с ffmpeg
)

func current_time() string { //получаем текущее время
	return time.Now().Format("15:04:05") //возвращаем время в формате HH:MM:SS
}

func current_date() string { //получаем текущую дату
	return time.Now().Format("02.01.2006") //возвращаем дату в DD.MM.YY
}

func file_directory(directory string) string { //получаем название директpeории, создаем, если не существует (по моему, выглядит как говнокод, когда возвращает входное значение при удачном результате)
	_, err := os.Stat(directory) //кидаем в _ информацию о директории
	if os.IsNotExist(err) {      //если нет директории, создаем
		os.Mkdir(directory, 0777) //создаем директорию с уровнем доступа 0777
	}
	return directory
}

func get_ffmpeg_config(config_directory string) []byte { //получаем JSON-файл с параметрами записи
	config_file, _ := os.ReadFile(config_directory) //читаем JSON файл
	return config_file
}

func cam_record(config_file []byte, file_directory string, file_name string) {
	var config Config                    //объявляем переменную config типа Config
	json.Unmarshal(config_file, &config) //анмаршалим JSON в структуру Config (TODO: прикутить обработчик ошибок)

	input := fmt.Sprintf("rtsp://%s:%s@%s/stream2", config.Cam_user, config.Cam_password, config.Cam_ip) //составляем параметр для input
	output := fmt.Sprintf("%s/%s.%s", file_directory, file_name, config.Filetype)                        //составляем параметр для output

	ffmpeg_go.Input(input, ffmpeg_go.KwArgs{"t": config.Duration}).Output(output, ffmpeg_go.KwArgs{"vcodec": "copy", "b:v": config.Bitrate}).Run() //запускаем ffmpeg

}

type Config struct { //структура JSON, аналогичная структуре файла
	Cam_ip       string `json:"cam_ip"`
	Cam_user     string `json:"cam_user"`
	Cam_password string `json:"cam_password"`
	Duration     string `json:"duration"`
	Filetype     string `json:"filetype"`
	Bitrate      string `json:"bitrate"`
}

func main() { //главная функция

	//config := get_ffmpeg_config(".config/config.json")
	//fmt.Printf("%s", config)

	//ffmpeg_go.Input("1.jpg").Output("1.png").Run()
	for {
		cfg := get_ffmpeg_config(".config/config.json")
		dir := file_directory(current_date())
		file := current_time()

		cam_record(cfg, dir, file)
	}
}

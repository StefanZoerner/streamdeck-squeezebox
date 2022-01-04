package squeezebox

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

const hostname = "elfman"
const port = 9090

func SetPlayerStatus(id string, status string) {

	connection_string := fmt.Sprintf("%s:%d", hostname, port)

	con, err := net.Dial("tcp", connection_string)
	checkErr(err)
	defer con.Close()

	cmd := fmt.Sprintf("%s mode %s\n", id, status);
	_, err = con.Write([]byte(cmd))
	checkErr(err)

	reply := make([]byte, 1024)
	_, err = con.Read(reply)
	checkErr(err)

	fmt.Println(string(reply))
}

func TogglePlayerStatus(id string) {

	connection_string := fmt.Sprintf("%s:%d", hostname, port)

	con, err := net.Dial("tcp", connection_string)
	checkErr(err)
	defer con.Close()

	cmd := fmt.Sprintf("%s mode ?\n", id);
	_, err = con.Write([]byte(cmd))
	checkErr(err)

	reply := make([]byte, 1024)
	_, err = con.Read(reply)
	checkErr(err)

	replyString := string(reply[:])
	replyString = strings.ReplaceAll(replyString, "\n", "");

	fmt.Printf("[%s]\n", replyString);

	parts := strings.Split(replyString, " ")
	current_mode := parts[2]

	fmt.Printf("current player mode: [%s]\n", current_mode);
	fmt.Println(strings.Count(current_mode, "play"));


	if strings.Count(current_mode, "play") == 1 {

		fmt.Println("play ->  pause")

		cmd := fmt.Sprintf("%s mode %s\n", id, "pause");
		_, err = con.Write([]byte(cmd))
		checkErr(err)

		reply = make([]byte, 1024)
		_, err = con.Read(reply)
		checkErr(err)

		fmt.Println(string(reply))
	} else if strings.Count(current_mode, "pause") == 1 || strings.Count(current_mode, "stop") == 1 {
		fmt.Println("pause/stop -> play")

		cmd := fmt.Sprintf("%s mode %s\n", id, "play");
		_, err = con.Write([]byte(cmd))
		checkErr(err)

		reply = make([]byte, 1024)
		_, err = con.Read(reply)
		checkErr(err)

		fmt.Println(string(reply))
	}
}

func CheckConnectionToPlayer(hostname string, port int, player_id string) error {

	connection_string := fmt.Sprintf("%s:%d", hostname, port)

	con, err := net.Dial("tcp", connection_string)
	if err != nil {
		return err
	}

	defer con.Close()
	cmd := fmt.Sprintf("%s connected ?\n", player_id);
	_, err = con.Write([]byte(cmd))
	if err != nil {
		return err
	}

	reply := make([]byte, 1024)
	_, err = con.Read(reply)
	if err != nil {
		return err
	}

	s_reply := string(reply)
	s_reply = strings.ReplaceAll(s_reply, "\n", "")
	result := strings.Split(s_reply," ")[2]

	if (strings.Contains(result, "%3F")) {
		return errors.New("Player "+player_id+" not connected to server.")
	}

	return nil
}

func GetCurrentArtworkUrl(hostname string, port int, player_id string) (string, error) {
	connection_string := fmt.Sprintf("%s:%d", hostname, port)

	con, err := net.Dial("tcp", connection_string)
	if err != nil {
		return "", err
	}

	defer con.Close()
	cmd := fmt.Sprintf("%s status - 1 tags:K,c\n", player_id);
	_, err = con.Write([]byte(cmd))
	if err != nil {
		return "", err
	}

	reply := make([]byte, 1024)
	_, err = con.Read(reply)
	if err != nil {
		return "", err
	}

	s_reply := string(reply)
	return s_reply, nil
}


func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}


func main () {
	url, err := GetCurrentArtworkUrl("elfman", 9090, "00:04:20:22:c2:54")
	if (err != nil) {
		fmt.Println("Fehler: " + err.Error())
	} else {
		fmt.Println("URL: " + url)
	}
}

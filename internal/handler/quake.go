package handler

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/maaarkin/quakecrawler/internal/domain"
	"github.com/maaarkin/quakecrawler/pkg/stringutil"
)

const (
	gameSuffix = "game_"
)

type quakeHandler struct {
}

func NewQuakeHandler() *quakeHandler {
	return &quakeHandler{}
}

func (q *quakeHandler) Run(scanner *bufio.Scanner) (map[string]domain.Payload, map[string]map[string]uint64) {
	isGameOpen := false
	gameCounter := 0
	report := make(map[string]domain.Payload)
	killByMeansReport := make(map[string]map[string]uint64)

	payload := domain.NewPayload()
	payloadKillByMeans := make(map[string]uint64)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "InitGame:") {
			if isGameOpen {
				gameCounter++
				gameId := fmt.Sprint(gameSuffix, gameCounter)
				report[gameId] = payload
				if len(payloadKillByMeans) > 0 {
					killByMeansReport[gameId] = payloadKillByMeans
				}
				payload = domain.NewPayload()
				payloadKillByMeans = make(map[string]uint64)
			}

			isGameOpen = true
		}

		if strings.Contains(scanner.Text(), "ClientUserinfoChanged:") {
			playerName := q.getPlayerName(scanner.Text())
			payload.AddPlayer(playerName)
		}

		if strings.Contains(scanner.Text(), "Kill:") {
			killer, death, deathType := q.getDataFromKill(scanner.Text())
			payload.AddKill(killer, death)
			payloadKillByMeans[deathType] = payloadKillByMeans[deathType] + 1
		}
	}

	if isGameOpen {
		gameCounter++
		gameId := fmt.Sprint(gameSuffix, gameCounter)
		report[gameId] = payload
		if len(payloadKillByMeans) > 0 {
			killByMeansReport[gameId] = payloadKillByMeans
		}
	}

	return report, killByMeansReport
}

func (q *quakeHandler) getPlayerName(line string) string {
	if !strings.Contains(line, "ClientUserinfoChanged:") {
		return ""
	}
	return stringutil.Between(line, `n\`, `\t`)
}

func (q *quakeHandler) getDataFromKill(line string) (killer, death, deathType string) {
	if !strings.Contains(line, "Kill:") {
		return "", "", ""
	}

	splitData := strings.Split(line, "Kill:")
	splitData = strings.Split(splitData[1], ":")

	description := splitData[1]

	killResult := strings.Split(description, "killed")
	killer = strings.Trim(killResult[0], " ")

	deathResult := strings.Split(killResult[1], "by")

	death = strings.Trim(deathResult[0], " ")
	deathType = strings.Trim(deathResult[1], " ")
	return
}

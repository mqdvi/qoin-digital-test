package main

import (
	"fmt"
	"math/rand"
)

type Player struct {
	ID          int
	Score       int
	Dice        []int
	RemovedDice []int
	Active      bool
	Next        int
}

func NewPlayers(numDice, numPlayers int) map[int]Player {
	players := make(map[int]Player)

	for i := 0; i < numPlayers; i++ {
		next := (i + 1) % numPlayers
		players[i] = Player{
			ID:     i,
			Score:  0,
			Dice:   make([]int, numDice),
			Active: true,
			Next:   next,
		}
	}

	return players
}

func (p *Player) RollDice() {
	for i := range p.Dice {
		p.Dice[i] = rand.Intn(6) + 1
	}
}

func (p *Player) Evaluate(players map[int]Player) {
	var newDice []int
	for _, die := range p.Dice {
		if die == 6 {
			p.Score++
		} else if die == 1 {
			p.RemovedDice = append(p.RemovedDice, 1)
		} else {
			newDice = append(newDice, die)
		}
	}

	var previous Player
	if p.ID == 0 {
		// if this is the first player
		// then assign previous to the last player
		previous = players[len(players)-1]
	} else {
		previous = players[p.ID-1]
	}

	if len(previous.RemovedDice) > 0 {
		newDice = append(newDice, previous.RemovedDice...)
		previous.RemovedDice = make([]int, 0)
		// update the previous player value in maps
		players[previous.ID] = previous
	}

	p.Dice = newDice
}

func (p *Player) Display() {
	if len(p.Dice) > 0 {
		fmt.Printf("\tPemain #%d (%d): %v\n", p.ID+1, p.Score, p.Dice)
	} else {
		fmt.Printf("\tPemain #%d (%d): _ (Berhenti bermain karena tidak memiliki dadu)\n", p.ID+1, p.Score)
	}
}

func main() {
	var (
		numPlayers,
		numDice int
	)

	fmt.Print("Masukkan jumlah player: ")
	fmt.Scan(&numPlayers)
	fmt.Print("Masukkan jumlah dadu: ")
	fmt.Scan(&numDice)

	if numPlayers <= 1 || numDice <= 0 {
		fmt.Println("Invalid input. Jumlah pemain minimal 2, dan jumlah dadu harus lebih dari 0.")
		return
	}

	players := NewPlayers(numDice, numPlayers)
	activePlayers := numPlayers
	round := 1

	fmt.Printf("\nPemain = %d, Dadu = %d\n", numPlayers, numDice)

	for activePlayers > 1 {
		fmt.Println("==================")
		fmt.Printf("Giliran %d lempar dadu:\n", round)

		for i, player := range players {
			player.RollDice()
			player.Display()
			players[i] = player
		}

		fmt.Println("Setelah evaluasi:")

		for i, player := range players {
			if player.Active {
				player.Evaluate(players)
			}

			player.Display()

			if len(player.Dice) == 0 && player.Active {
				player.Active = false
				activePlayers--
			}

			players[i] = player
		}

		round++
	}

	lastPlayer, winningPlayer, draw := findWinner(players, numPlayers)

	fmt.Println("==================")
	fmt.Printf("Game berakhir karena hanya pemain #%d yang memiliki dadu.\n", lastPlayer.ID+1)

	if draw {
		fmt.Println("Game berakhir seri karena pemain memiliki point tertinggi yang sama.")
	} else {
		fmt.Printf("Game dimenangkan oleh pemain #%d karena memiliki poin lebih banyak dari pemain lainnya.\n", winningPlayer.ID+1)
	}
}

func findWinner(players map[int]Player, numPlayers int) (Player, Player, bool) {
	var lastPlayer, winningPlayer Player
	draw := false
	highestScore := 0

	for _, p := range players {
		if p.Active {
			lastPlayer = p
			continue
		}

		if p.Score > highestScore {
			winningPlayer = p
			highestScore = p.Score
		} else if p.Score == highestScore {
			draw = true
		}
	}

	return lastPlayer, winningPlayer, draw
}

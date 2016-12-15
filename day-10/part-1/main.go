package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
}

const ANON int = -999

func main() {
	s := bufio.NewScanner(os.Stdin)
	m := NewModel()

	for s.Scan() {
		t := strings.Split(s.Text(), " ")
		if len(t) < 5 {
			panic(fmt.Sprintf("Unexpectedly short input: %q", s.Text()))
		}

		if t[0] == "value" {
			if len(t) != 6 {
				panic(fmt.Sprintf("Expected input of 6 words but got %q", s.Text()))
			}

			value, _ := strconv.Atoi(t[1])
			bot, _ := strconv.Atoi(t[5])

			m.GiveValue(value, bot, ANON)
			continue
		}

		if t[0] == "bot" {
			if len(t) != 12 {
				panic(fmt.Sprintf("Expected input of 12 words but got %q", s.Text()))
			}

			fromBot, _ := strconv.Atoi(t[1])
			lowBot, _ := strconv.Atoi(t[6])
			highBot, _ := strconv.Atoi(t[11])

			if t[5] == "output" {
				lowBot = ^lowBot
			}
			if t[10] == "output" {
				highBot = ^highBot
			}

			m.AssignRecipient(fromBot, lowBot, highBot)
			continue
		}

		panic(fmt.Sprintf("Unrecognised input: %q", s.Text()))
	}

	if err := s.Err(); err != nil {
		panic(err)
	}
}

type Bot struct {
	id int

	low, high int
	holding   int

	botLow, botHigh int
	recipientsSet   bool
}

func (b *Bot) String() string {
	if b.id == ANON {
		return "Input"
	}
	if b.id < 0 {
		return fmt.Sprintf("Output %d", ^b.id)
	}
	return fmt.Sprintf("Bot %d", b.id)
}

type Model struct {
	bots map[int]*Bot
}

func NewModel() *Model {
	return &Model{
		bots: make(map[int]*Bot),
	}
}

func (m *Model) bot(n int) *Bot {
	bot, found := m.bots[n]
	if found {
		return bot
	}

	m.bots[n] = &Bot{id: n}
	log.Printf("Encountered %s.", m.bots[n])
	return m.bots[n]
}

func (m *Model) GiveValue(value, toBot, fromBot int) {
	bot := m.bot(toBot)
	log.Printf("%s received value %d from %s.", bot, value, m.bot(fromBot))

	if bot.holding == 0 {
		log.Printf("%s holding on to value %d.", bot, value)
		bot.low = value
		bot.holding = 1
		return
	}

	bot.low, bot.high = lohi(bot.low, value)
	bot.holding = 2

	if !bot.recipientsSet {
		log.Printf("%s holding low value %d and high value %d until recipients set.", bot, bot.low, bot.high)
		return
	}
	m.forward(bot)
}

func (m *Model) forward(bot *Bot) {
	log.Printf("%s forwarding low value %d to %s; and high value %d to %s.", bot, bot.low, m.bot(bot.botLow), bot.high, m.bot(bot.botHigh))

	m.GiveValue(bot.low, bot.botLow, bot.id)
	m.GiveValue(bot.high, bot.botHigh, bot.id)

	bot.low, bot.high, bot.holding = 0, 0, 0
}

func (m *Model) AssignRecipient(fromBot, lowBot, highBot int) {
	from, low, high := m.bot(fromBot), m.bot(lowBot), m.bot(highBot)
	log.Printf("%s will forward low value to %s; and high value to %s.", from, low, high)

	from.botLow = low.id
	from.botHigh = high.id
	from.recipientsSet = true

	if from.holding == 2 {
		log.Printf("%s was awaiting recipients. Forwarding now.", from)
		m.forward(from)
	}
}

func lohi(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}

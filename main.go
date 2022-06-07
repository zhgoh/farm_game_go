package main

import (
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var seed_type = []string{"carrot", "corn", "ginger", "onion", "tomato", "wheat"}
var seed_price = []float32{2.0, 3.0, 4.5, 6.0, 7.0, 9.0}
var seed_durations = []int{1, 2, 3, 4, 5, 6}
var crop_price = []float32{2.5, 4.0, 6.0, 11.0, 14.0, 20.0}

const time_delay = 5

func main() {
	current_farm := init_farm()
	go func() {
		for {
			time.Sleep(time_delay * time.Second)
			for i := 0; i < len(current_farm.plants); i++ {
				if current_farm.plants[i].priority > 0 {
					item := current_farm.plants[i]
					current_farm.plants.update(item, item.value, item.priority-1)
				}
			}
		}
	}()
	intro(current_farm)
}

func intro(current_farm *farm) {
	fmt.Println("Welcome to your farm.")
	fmt.Printf("What do you want to name your farm? ")
	showMainOptions(current_farm)
}

func clearScreen() {
	fmt.Printf("\033[H\033[2J")
}

func showStats(current_farm *farm) {
	fmt.Printf("Money: $%.2f, Lands: %d \n",
		current_farm.money,
		current_farm.land_free,
	)

	fmt.Println()
	fmt.Println("Seeds:")
	for k, v := range seed_type {
		fmt.Printf("%s: %d, ", v, current_farm.seeds[k])
	}

	fmt.Println()
	fmt.Println()

	// fmt.Printf("%+v\n", current_farm.seeds)
	// for _, v := range current_farm.plants {
	// 	fmt.Printf("%d %f, ", v.priority, v.value)
	// }
	// fmt.Println()
}

func showMainOptions(current_farm *farm) {
	for {
		clearScreen()
		showStats(current_farm)

		fmt.Println("Type one of the options below.")
		fmt.Println()

		fmt.Println("(p)lant seeds.            >")
		fmt.Println("(h)arvest and sell crops. >")
		fmt.Println("(b)uy seeds/land.         >")
		fmt.Println("(s)ave game.")
		fmt.Println("(l)oad game.")
		fmt.Println("(q)uit game.")

		fmt.Println()
		fmt.Printf("#> ")

		var option string
		fmt.Scanf("%s\n", &option)

		switch option {
		case "p":
			showPlantSeedOptions(current_farm)
		case "h":
			showHarvestOptions(current_farm)
		case "b":
			showBuySeedsOptions(current_farm)
		case "s":
			// Save game
			file, err := os.Create("game.txt")
			defer file.Close()

			if err != nil {
				fmt.Println(err)
			}

			fmt.Fprintf(file,
				"Money: %.2f\nLand: %d\nSeeds: %d\n",
				current_farm.money,
				current_farm.land_free,
				len(current_farm.seeds))

			for i := range current_farm.seeds {
				fmt.Fprintf(file, "%d ", current_farm.seeds[i])
			}
			fmt.Fprintln(file)
			fmt.Fprintf(file, "plants: %d\n", current_farm.plants.Len())
			for _, v := range current_farm.plants {
				fmt.Fprintf(file, "%v", v)
			}

			fmt.Println("Saved game.")
			fmt.Println("<Press Enter to continue>")
			fmt.Scanln()

		case "l":
			// Load game
			file, err := os.Open("game.txt")
			defer file.Close()

			if err != nil {
				fmt.Println(err)
			}

			var money float32
			var land_free int
			var numSeeds int
			fmt.Fscanf(file,
				"Money: %f\nLand: %d\nSeeds: %d\n",
				&money,
				&land_free,
				&numSeeds)

			seeds := make([]int, numSeeds)
			for i := 0; i < numSeeds; i++ {
				var val int
				_, err := fmt.Fscanf(file, "%d", &val)
				if err != nil {
					fmt.Println("Error reading values from seed")
				}
				seeds[i] = val
			}
			current_farm.money = money
			current_farm.land_free = land_free
			current_farm.seeds = seeds

			plants := 0
			fmt.Fscanf(file, "plants: %d\n", &plants)
			current_farm.plants = make(PriorityQueue, plants)
			for i := 0; i < plants; i++ {
				item := &Item{}
				fmt.Fscanf(file, "%v", &item)
				current_farm.plants = append(current_farm.plants, item)
			}

			fmt.Println("Loaded game.")
			fmt.Println("<Press Enter to continue>")
			fmt.Scanln()

		case "q":
			fmt.Println("Quitting the game.")
			return
		default:
			fmt.Println("No such option.")
			fmt.Println("<Press Enter to continue>")
			fmt.Scanln()
		}
	}
}

func showPlantSeedOptions(current_farm *farm) {
	for {
		clearScreen()
		showStats(current_farm)

		fmt.Println("Type one of the options below.")
		fmt.Println()

		for i := 0; i < len(seed_type); i++ {
			fmt.Printf("(%d) plant %s (%d), time: %d, price: $%.2f.\n", i, seed_type[i], current_farm.seeds[i], seed_durations[i], seed_price[i])
		}
		fmt.Println("(q)uit to previous menu.")
		fmt.Println()
		fmt.Printf("#> ")

		var option string
		fmt.Scanf("%s\n", &option)

		if num, ok := strconv.Atoi(option); ok == nil {
			if num >= 0 && num < len(seed_type) {
				// Check if enough seeds and loads
				if current_farm.seeds[num] == 0 {
					fmt.Printf("Not enough %[1]s seeds. Buy more %[1]s seeds to plant.\n", seed_type[num])
					fmt.Println("<Press Enter to continue>")
					fmt.Scanln()
				} else if current_farm.land_free == 0 {
					fmt.Println("Not enough lands to grow. Buy more lands to plant.")
					fmt.Println("<Press Enter to continue>")
					fmt.Scanln()
				} else {
					fmt.Printf("Planting %s\n", seed_type[num])
					fmt.Println("<Press Enter to continue>")
					fmt.Scanln()
					current_farm.seeds[num] -= 1

					item := &Item{
						value:    crop_price[num],
						priority: seed_durations[num],
					}
					heap.Push(&current_farm.plants, item)
					current_farm.land_free -= 1
				}
			} else {
				fmt.Println("No such seed. ")
				fmt.Println("<Press Enter to continue>")
				fmt.Scanln()
			}
		} else if strings.ToLower(option) == "q" {
			return
		} else {
			fmt.Println("No such option. ")
			fmt.Println("<Press Enter to continue>")
			fmt.Scanln()
		}
	}
}

func showHarvestOptions(current_farm *farm) {
	for {
		clearScreen()
		showStats(current_farm)

		fmt.Println("Type one of the options below.")
		fmt.Println()

		fmt.Println("(h)arvest and sell all crops.")
		fmt.Println("(q)uit to previous menu.")
		fmt.Println()
		fmt.Printf("#> ")

		var option string
		fmt.Scanf("%s\n", &option)

		// Intentionally left out checks for other keys so that we can spam enter key on this page
		switch strings.ToLower(option) {
		case "h":
			count := 0
			for current_farm.plants.Len() > 0 {
				if current_farm.plants[0].priority > 0 {
					break
				}

				item := heap.Pop(&current_farm.plants).(*Item)
				current_farm.land_free += 1

				current_farm.money += item.value
				count += 1
			}
			fmt.Printf("Harvested and sold %d crop(s).\n", count)
			fmt.Println("<Press Enter to continue>")
			fmt.Scanln()
		case "q":
			return
		}
	}
}

func showBuySeedsOptions(current_farm *farm) {
	for {
		clearScreen()
		showStats(current_farm)
		fmt.Println("Type one of the options below.")
		fmt.Println()

		k, v := 0, ""
		for k, v = range seed_type {
			fmt.Printf("(%d) buy %s seed - $%.2f.\n", k, v, seed_price[k])
		}

		fmt.Printf("(%d) buy land - $50.00.\n", k+1)
		fmt.Println("(q)uit to previous menu.")
		fmt.Println()
		fmt.Printf("#> ")

		var option string
		fmt.Scanf("%s\n", &option)

		switch strings.ToLower(option) {
		case "q":
			return
		default:
			if val, err := strconv.Atoi(option); err == nil {
				if val == k+1 {
					if current_farm.money > 50 {
						fmt.Println("Buying 1 land.")
						current_farm.money -= 50
						current_farm.land_free += 1
					} else {
						fmt.Println("Not enough money to buy land.")
					}
				} else {
					if val >= 0 && val < len(seed_type) {
						if current_farm.money >= seed_price[val] {
							fmt.Printf("Buying %s seed.\n", seed_type[val])
							current_farm.money -= seed_price[val]
							current_farm.seeds[val]++
						} else {
							fmt.Println("Not enough money to buy seed.")

						}
					} else {
						fmt.Println("No such seed found.")
					}
				}

				fmt.Println("<Press Enter to continue>")
				fmt.Scanln()
			}
		}
	}
}

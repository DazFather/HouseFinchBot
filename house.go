package main

import "time"

type House struct {
	Chart   ShoppingList
	ownerID int64
	shared  map[int64]string
	active  chan bool
	cleaner *time.Timer
}

var CACHE = make(map[int64]*House)

func newHouse(ownerID int64) (house *House) {
	// Consider an house unsued after no given interaction in:
	const UNUSED = time.Hour * 24 * 240
	// Create the channel for each user interaction
	var interact = make(chan bool)

	// Initialize the house
	house = &House{
		Chart:   newShoppingList(),
		ownerID: ownerID,
		active:  interact,
		cleaner: time.AfterFunc(UNUSED, func() {
			interact <- false
		}),
	}
	// save it
	CACHE[ownerID] = house

	// In a background process...
	go func() {
		// ... check if user are not interacting ...
		for active := range interact {
			if !active || !house.cleaner.Stop() {
				break
			}
			house.cleaner.Reset(UNUSED)
		}
		// ... in this case delete the house from memory
		house.delete()
	}()

	return
}

func (house *House) delete() {
	close(house.active)
	for userID := range house.shared {
		delete(CACHE, userID)
	}
	delete(CACHE, house.ownerID)
}

func SelectHouse(chatID int64) (house *House, registeredUser bool) {
	house, registeredUser = CACHE[chatID]
	if !registeredUser {
		house = newHouse(chatID)
	}
	house.active <- true
	return
}

func (house House) IsOwner(userID int64) bool {
	return house.ownerID == userID
}

func (house House) IsMember(userID int64) bool {
	return CACHE[house.ownerID] == CACHE[userID]
}

func (house *House) Join(userID int64, name string) bool {
	if house.IsOwner(userID) || house.IsMember(userID) {
		return false
	}

	CACHE[userID] = house
	if len(house.shared) == 0 {
		house.shared = make(map[int64]string)
	}
	house.shared[userID] = name
	return true
}

func (house House) Members() int {
	return len(house.shared)
}

func (house House) EachMembers(do func(userID int64, name string)) {
	for userID, name := range house.shared {
		do(userID, name)
	}
}

func (house *House) Kick(userID int64) bool {
	if house.IsOwner(userID) {
		return false
	}

	for id := range house.shared {
		if id == userID {
			delete(house.shared, userID)
			delete(CACHE, userID)
			return true
		}
	}
	return false
}

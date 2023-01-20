package main

type House struct {
	Chart   ShoppingList
	ownerID int64
	shared  map[int64]string
}

func newHouse(ownerID int64) *House {
	return &House{Chart: newShoppingList(), ownerID: ownerID}
}

func SelectHouse(chatID int64) (house *House, registeredUser bool) {
	house, registeredUser = CACHE[chatID]
	if !registeredUser {
		house = newHouse(chatID)
		CACHE[chatID] = house
	}
	return
}

var CACHE = make(map[int64]*House)

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

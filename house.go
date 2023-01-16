package main

type House struct {
	Chart   ShoppingList
	ownerID int64
	shared  []int64
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

func (house *House) Join(userID int64) bool {
	if house.IsOwner(userID) || house.IsMember(userID) {
		return false
	}

	CACHE[userID] = house
	if len(house.shared) == 0 {
		house.shared = []int64{userID}
		return true
	}
	house.shared = append(house.shared, userID)
	return true
}

func (house *House) Kick(userID int64) bool {
	if house.IsOwner(userID) {
		return false
	}

	for ind, id := range house.shared {
		if id == userID {
			house.shared = append(house.shared[:ind], house.shared[ind+1:]...)
			delete(CACHE, userID)
			return true
		}
	}
	return false
}

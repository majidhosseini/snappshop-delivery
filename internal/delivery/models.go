package delivery

type Request struct {
    OrderID          string   `json:"orderId"`
    UserInfo         UserInfo `json:"userInfo"`
    FromLoc          [2]float64 `json:"fromLoc"`
    ToLoc            [2]float64 `json:"toLoc"`
    DeliveryTimeFrame [2]string `json:"deliveryTimeFrame"`
    State            string     `json:"state"`
}

type UserInfo struct {
    Name    string `json:"name"`
    Phone   string `json:"phone"`
    Address string `json:"address"`
}

type DeliveryState string

const (
    StateInit       DeliveryState = "init"
    StateIsFinding  DeliveryState = "isFinding"
    StateFound      DeliveryState = "found"
    StateNotFound   DeliveryState = "notFound"
    StateDelivered  DeliveryState = "delivered"
)

package jcdportal

// Station reflects the location, information and status of a JCDecaux Bike Renting facility
type Station struct {
	Name         string `json:"name"`         // Name of the station
	ContractName string `json:"contractName"` // Name of the contract of the station

	// Number of the station. This is NOT an id, thus it is unique only inside a contract.
	Number int `json:"number"`

	// Address of the station. As it is raw data, sometimes it will be more of a comment than an address.
	Address string `json:"address"`

	Position Position `json:"position"` // Position of the station in WGS84 format

	Status     string `json:"status"`     // Status indicates whether this station is 'CLOSED' or 'OPEN'
	LastUpdate string `json:"lastUpdate"` // LastUpdate is a timestamp indicating the last update time

	Banking   bool `json:"banking"`   // Banking indicates whether this station has a payment terminal
	Bonus     bool `json:"bonus"`     // Bonus indicates whether this is a bonus station
	Connected bool `json:"connected"` // Connected indicates whether the station is connected to it's backend
	Overflow  bool `json:"overflow"`  // Overflow indicates if the station allows overflow bike return

	// TotalStands indicates the total bike capacity, the number of empty places and the total number of bikes on the station
	TotalStands Stands `json:"totalStands"`

	// MainStands indicates the physical bike stand capacity, the number of empty bike stands and the number bikes locked on stands
	MainStands Stands `json:"mainStands"`

	// OverflowStands indicates the overflow bike capacity, the number of empty overflow places and the number of bikes in overflow
	OverflowStands Stands `json:"overflowStands"`
}

// Park reflects the location, information and status of a JCDecaux Bike Parking facility
type Park struct {
	Name         string `json:"name"`         // Name of the park
	Number       int    `json:"number"`       // The identification number of the park
	ContractName string `json:"contractName"` // Contract name of the park

	// Address of the park. As it is raw data, sometimes it will be more of a comment than an address.
	Address string `json:"address"`
	ZipCode string `json:"zipCode"` // Park's zip code

	Status   string   `json:"status"`   // Status indicates whether this station is 'CLOSED' or 'OPEN'
	Position Position `json:"position"` // Position of the station in WGS84 format

	AccessType string `json:"accessType"` // 'FREE_ACCESS' or 'SECURED' indicates of the park access is free or secured
	LockerType string `json:"lockerType"` // 'SINGLE' or 'COLLECTIVE' indicates whether it's a collective or single park
	City       string `json:"city"`       // Park's city

	HasSurveillance      bool `json:"hasSurveillance"`      // HasSurveillance indicates whether the park is guarded
	IsFree               bool `json:"isFree"`               // IsFree indicates if the park is free of charge
	IsOffStreet          bool `json:"isOffStreet"`          // IsOffStreet indicates whether it's off street
	HasElectricSupport   bool `json:"hasElectricSupport"`   // HasElectricSupport indicates whether the park handle eBikes charge
	HasPhysicalReception bool `json:"hasPhysicalReception"` // HasPhysicalReception indicates whether the park has a customer office
}

// Contract reflects a given bike contract that JCDecaux is currently operating
type Contract struct {
	Name        string   `json:"name"`         // is the identifier of the contract
	CountryCode string   `json:"country_code"` // CountryCode is the code (ISO 3166) of the country
	Cities      []string `json:"cities"`       // Cities that are concerned by this contract

	// CommercialName is the commercial name of the contract (the one users usually know)
	CommercialName string `json:"commercial_name"`
}

// Stands indicates bike stand state
type Stands struct {
	Capacity       int            `json:"capacity"`
	Availabilities Availabilities `json:"availabilities"`
}

type Availabilities struct {
	Bikes  int `json:"bikes"`
	Stands int `json:"stands"`
}

// Position of the station in WGS84 format
type Position struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

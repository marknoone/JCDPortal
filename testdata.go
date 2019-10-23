package jcdportal

var (
	DummyStation = Station{
		Number:       123,
		ContractName: "Lyon",
		Name:         "nom station",
		Address:      "adresse indicative",
		Position: Position{
			Latitude:  45.774204,
			Longitude: 4.867512,
		},
		Banking:    true,
		Bonus:      false,
		Status:     "OPEN",
		LastUpdate: "2019-04-08T12:23:34Z",
		Connected:  true,
		Overflow:   true,
		TotalStands: Stands{
			Availabilities: Availabilities{
				Bikes:  15,
				Stands: 25,
			},
			Capacity: 40,
		},
		MainStands: Stands{
			Availabilities: Availabilities{
				Bikes:  15,
				Stands: 15,
			},
			Capacity: 30,
		},
		OverflowStands: Stands{
			Availabilities: Availabilities{
				Bikes:  0,
				Stands: 10,
			},
			Capacity: 10,
		},
	}

	DummyPark = Park{
		ContractName: "nantes",
		Name:         "GARE DE PONT ROUSSEAU NORD",
		Number:       89,
		Status:       "OPEN",
		Position: Position{
			Latitude:  47.1920011,
			Longitude: -1.5490259,
		},
		AccessType:           "FREE_ACCESS",
		LockerType:           "SINGLE",
		HasSurveillance:      false,
		IsFree:               true,
		Address:              "Rue de la Gare",
		ZipCode:              "44400",
		City:                 "Rezé",
		IsOffStreet:          true,
		HasElectricSupport:   false,
		HasPhysicalReception: false,
	}

	DummyContract = Contract{
		Name:           "Lyon",
		CommercialName: "Vélo'v",
		CountryCode:    "FR",
		Cities: []string{
			"Lyon",
			"Villeurbanne",
		},
	}
)

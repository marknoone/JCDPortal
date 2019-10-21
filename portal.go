package jcdportal

type BikesAPI struct {
	Contract string
	apiKey   string
}

func (b BikesAPI) request() {}

func (b BikesAPI) GetContracts() {}

func (b BikesAPI) GetStations() {}

func (b BikesAPI) GetContractStations() {}

func (b BikesAPI) GetContractParks() {}

func (b BikesAPI) GetContract() {}

func (b BikesAPI) GetStation() {}

func (b BikesAPI) GetPark() {}

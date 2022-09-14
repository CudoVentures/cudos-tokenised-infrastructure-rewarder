package tokenised_infrastructure_rewarder

import (
	"github.com/CudoVentures/tokenised-infrastructure-rewarder/internal/app/tokenised-infrastructure-rewarder/infrastructure"
	"github.com/CudoVentures/tokenised-infrastructure-rewarder/internal/app/tokenised-infrastructure-rewarder/services"
	"github.com/CudoVentures/tokenised-infrastructure-rewarder/internal/app/tokenised-infrastructure-rewarder/types"
	"github.com/rs/zerolog/log"
)

func Start() error {

	farms := getDummyData() // replace with call to backend once it is done
	log.Info().Msgf("Farms fetched from backend: %s", farms)
	config := infrastructure.NewConfig()
	err := services.ProcessPayment(config)
	return err
	// return nil
}

func getDummyData() types.Farm {

	Collection := types.Collection{Denom: types.Denom{Id: "test"}, Nfts: []types.NFT{}}
	testFarm := types.Farm{Id: "test", SubAccountName: "test", BTCWallet: "testwallet2", Collections: []types.Collection{Collection}}
	return testFarm
}

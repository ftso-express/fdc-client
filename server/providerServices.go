package server

import (
	"encoding/hex"
	"flare-common/restServer"
	"fmt"
	"math/rand"

	"github.com/ethereum/go-ethereum/crypto"
)

type addData struct {
	data       string
	additional string
}

// TODO: Luka - Get this from config
const PROVIDER_RANDOM_SEED = 42

func calculateMaskedRoot(real_root string, random_num string) string {
	return hex.EncodeToString(crypto.Keccak256([]byte(real_root), []byte(random_num)))
}

func (controller *FDCProtocolProviderController) saveRoot(address string, round uint64, root string, random string) {
	if restServer.IsNil(controller.rootStorage) {
		controller.rootStorage = make(map[string]map[uint64]merkleRootStorageObject)
	}
	if _, ok := controller.rootStorage[address]; !ok {
		controller.rootStorage[address] = make(map[uint64]merkleRootStorageObject)
	}
	controller.rootStorage[address][round] = merkleRootStorageObject{merkleRoot: root, randomNum: random}
}

func (controller *FDCProtocolProviderController) submit1Service(round uint64, address string) (string, error) {
	fmt.Println("Submit1Handler")
	fmt.Printf("round: %s\n", fmt.Sprint(round))
	fmt.Printf("address: %s\n", address)
	return fmt.Sprint(round) + address, nil
}

func (controller *FDCProtocolProviderController) submit2Service(round uint64, address string) (string, error) {
	// Get merkle tree root from attestation client from controller
	r1 := rand.New(rand.NewSource(int64(round)))
	real_root := hex.EncodeToString(crypto.Keccak256([]byte(fmt.Sprintf("%X", r1.Int63()))))

	r2 := rand.New(rand.NewSource(PROVIDER_RANDOM_SEED))
	random_num := hex.EncodeToString(crypto.Keccak256([]byte(fmt.Sprintf("%X", r2.Int63()))))

	// save root to storage
	controller.saveRoot(address, round, real_root, random_num)

	masked := calculateMaskedRoot(real_root, random_num)

	return masked, nil
}

func (controller *FDCProtocolProviderController) submitSignaturesService(round uint64, address string) (addData, error) {
	// check storage if root was saved
	if _, ok := controller.rootStorage[address]; !ok {
		return addData{}, fmt.Errorf("address not in storage")
	}
	if _, ok := controller.rootStorage[address][round]; !ok {
		return addData{}, fmt.Errorf("round for address not in storage")
	}
	savedRoot := controller.rootStorage[address][round]

	fmt.Println("SubmitSignaturesHandler")
	fmt.Printf("round: %s\n", fmt.Sprint(round))
	fmt.Printf("address: %s\n", address)
	fmt.Printf("root: %s\n", savedRoot.merkleRoot)
	fmt.Printf("random: %s\n", savedRoot.randomNum)

	return addData{data: savedRoot.merkleRoot, additional: savedRoot.randomNum}, nil
}

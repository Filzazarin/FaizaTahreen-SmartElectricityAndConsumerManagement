#!/bin/bash
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error
set -e

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1
starttime=$(date +%s)
CC_SRC_LANGUAGE=${1:-"go"}
CC_SRC_LANGUAGE=`echo "$CC_SRC_LANGUAGE" | tr [:upper:] [:lower:]`
if [ "$CC_SRC_LANGUAGE" = "go" -o "$CC_SRC_LANGUAGE" = "golang"  ]; then
	CC_RUNTIME_LANGUAGE=golang
	CC_SRC_PATH=github.com/chaincode/electricity/go
elif [ "$CC_SRC_LANGUAGE" = "go1" ]; then
	CC_RUNTIME_LANGUAGE=golang
	CC_SRC_PATH=github.com/chaincode/electricity/go1
elif [ "$CC_SRC_LANGUAGE" = "go2" ]; then
	CC_RUNTIME_LANGUAGE=golang
	CC_SRC_PATH=github.com/chaincode/electricity/go2
else
	echo The chaincode folder ${CC_SRC_LANGUAGE} is not available by this script
	echo Supported chaincode languages are: go
	exit 1
fi


# clean the keystore
rm -rf ./hfc-key-store

# launch network; create channel and join peer to channel
cd ../first-network
echo y | ./byfn.sh down
echo y | ./byfn.sh up -a -n -s couchdb

CONFIG_ROOT=/opt/gopath/src/github.com/hyperledger/fabric/peer
ORG1_MSPCONFIGPATH=${CONFIG_ROOT}/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
ORG1_TLS_ROOTCERT_FILE=${CONFIG_ROOT}/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
ORG2_MSPCONFIGPATH=${CONFIG_ROOT}/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
ORG2_TLS_ROOTCERT_FILE=${CONFIG_ROOT}/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
ORDERER_TLS_ROOTCERT_FILE=${CONFIG_ROOT}/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
set -x

echo "Installing smart contract on peer0.org1.example.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=Org1MSP \
  -e CORE_PEER_ADDRESS=peer0.org1.example.com:7051 \
  -e CORE_PEER_MSPCONFIGPATH=${ORG1_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${ORG1_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n electricity \
    -v 1.0 \
    -p "$CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"

echo "Installing smart contract on peer0.org2.example.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=Org2MSP \
  -e CORE_PEER_ADDRESS=peer0.org2.example.com:9051 \
  -e CORE_PEER_MSPCONFIGPATH=${ORG2_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${ORG2_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n electricity \
    -v 1.0 \
    -p "$CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"

echo "Instantiating smart contract on mychannel"
docker exec \
  -e CORE_PEER_LOCALMSPID=Org1MSP \
  -e CORE_PEER_MSPCONFIGPATH=${ORG1_MSPCONFIGPATH} \
  cli \
  peer chaincode instantiate \
    -o orderer.example.com:7050 \
    -C mychannel \
    -n electricity \
    -l "$CC_RUNTIME_LANGUAGE" \
    -v 1.0 \
    -c '{"Args":[]}' \
    -P "AND('Org1MSP.member','Org2MSP.member')" \
    --tls \
    --cafile ${ORDERER_TLS_ROOTCERT_FILE} \
    --peerAddresses peer0.org1.example.com:7051 \
    --tlsRootCertFiles ${ORG1_TLS_ROOTCERT_FILE}

echo "Waiting for instantiation request to be committed ..."
sleep 10

echo "Submitting initLedger transaction to smart contract on mychannel"
echo "The transaction is sent to the two peers with the chaincode installed (peer0.org1.example.com and peer0.org2.example.com) so that chaincode is built before receiving the following requests"
docker exec \
  -e CORE_PEER_LOCALMSPID=Org1MSP \
  -e CORE_PEER_MSPCONFIGPATH=${ORG1_MSPCONFIGPATH} \
  cli \
  peer chaincode invoke \
    -o orderer.example.com:7050 \
    -C mychannel \
    -n electricity \
    -c '{"function":"initLedger","Args":[]}' \
    --waitForEvent \
    --tls \
    --cafile ${ORDERER_TLS_ROOTCERT_FILE} \
    --peerAddresses peer0.org1.example.com:7051 \
    --peerAddresses peer0.org2.example.com:9051 \
    --tlsRootCertFiles ${ORG1_TLS_ROOTCERT_FILE} \
    --tlsRootCertFiles ${ORG2_TLS_ROOTCERT_FILE}
set +x

cat <<EOF

Total setup execution time : $(($(date +%s) - starttime)) secs ...

Next, use the Electricity applications to interact with the deployed Electricity chaincode.
The Electricity applications are available in JavaScript programming language using Node.js SDK.
Follow the instructions for the chaincode you installed in the network while running the shell script startFabric.sh:

Consumer details:

  Start by changing into the "ConsumerDetails" directory:
    cd ConsumerDetails

  Next, install all required packages(only for the first time):
    npm install

  Then run the following applications to enroll the admin user, and register a new user
  called user1 which will be used by the other applications to interact with the deployed
  ConsumerDetails chaincode:
    node enrollAdmin
    node registerUser

  You can run the invoke application as follows. By default, the invoke application will
  add a new connection detail for a consumer in the ledger, but you can update the application to submit other transactions:
    node invoke

  You can run the query application as follows. By default, the query application will
  return all consumer details stored in the ledger, but you can update the application to evaluate other transactions:
    node query

Bill Details:

  Start by changing into the "BillDetails" directory:
    cd BillDetails

  Next, install all required packages(only for the first time):
    npm install

  Then run the following applications to enroll the admin user, and register a new user
  called user1 which will be used by the other applications to interact with the deployed
  BillDetails chaincode:
    node enrollAdmin
    node registerUser

  You can run the invoke application as follows. By default, the invoke application will
  add a new bill detail for a consumer in the ledger, but you can update the application to submit other transactions:
    node invoke

  You can run the query application as follows. By default, the query application will
  return all consumer bill details stored in the ledger, but you can update the application to evaluate other transactions:
    node query
    
Payment details:

  Start by changing into the "PaymentDetails" directory:
    cd PaymentDetails

  Next, install all required packages(only for the first time):
    npm install

  Then run the following applications to enroll the admin user, and register a new user
  called user1 which will be used by the other applications to interact with the deployed
  PaymentDetails contract:
    node enrollAdmin
    node registerUser

  You can run the invoke application as follows. By default, the invoke application will
  add a new payment detail for a consumer in the ledger, but you can update the application to submit other transactions:
    node invoke

  You can run the query application as follows. By default, the query application will
  return all consumer payment details stored in the ledger, but you can update the application to evaluate other transactions:
    node query



EOF

DIR=`dirname "$0"`

# config chain id
bandoracled2 config chain-id bandchain

# add validator to bandoracled2 config
bandoracled2 config validator bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec

# send band tokens to reporter
echo "benefit select evidence crystal nature shell arm struggle sibling thank dish cruel immense erode coil inmate brave tackle gas rural giggle welcome next aspect"
bandcli tx send validator band1r8sht7m0veaqgx2zpqydxmwy9f6v3c995cxh7n 1000000uband --keyring-backend test

# wait for sending transaction success
sleep 2

# add reporter to bandchain
bandcli tx oracle add-reporter band1r8sht7m0veaqgx2zpqydxmwy9f6v3c995cxh7n --from validator --keyring-backend test

# run bandoracled2
bandoracled2 run
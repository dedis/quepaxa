pwd=$(pwd)
. "${pwd}"/experiments/setup-5/ip.sh

echo "Asynchrony tests"

python3 experiment/asynchrony/test_automation.py 5
python3 experiment/asynchrony/summary.py 5 ${device}
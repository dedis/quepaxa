pwd=$(pwd)
. "${pwd}"/experiments/setup-5/ip.sh

echo "Asynchrony tests"

python3 experiments/asynchrony/test_automation.py 15
python3 experiments/asynchrony/summary.py 15 ${device}
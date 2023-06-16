mode=$1

if [[ "${mode}" == "LAN" ]]
then
  echo "Running best-case LAN"
  # LAN experiment
  python3 experiment/best-case/test-automation.py LAN 5
  python3 experiment/best-case/summary.py LAN 5
fi

if [[ "${mode}" == "WAN" ]]
then
  echo "Running best-case WAN"
  # WAN experiment
  python3 experiment/best-case/test-automation.py WAN 5
  python3 experiment/best-case/summary.py WAN 5
  echo "Best case experiments done"
fi
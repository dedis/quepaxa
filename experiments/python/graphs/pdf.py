import sys
import matplotlib.pyplot as plt
import numpy as np
import scipy.stats as st

filename=sys.argv[1]

S = []

with open(filename) as f:
    content = f.readlines()
    content=content[1:]
    for row in content:
        decided = row.strip().split(" ")[0]
        step = row.strip().split(" ")[8]
        if(decided == "decided"):
            S.append(int(step))


def draw_histogram(latency_values):
    plt.hist(latency_values)
    plt.rcParams.update({'font.size': 13.45})
    plt.xlabel('Step (S)')
    plt.ylabel('Count')
    plt.legend()
    plt.show()
    plt.close()

def get_pdf(latency_list):
    np_array = np.array(latency_list)  # convert the list into a numpy array
    ag = st.gaussian_kde(np_array)  # calculate the kernel density function for the latency values
    # list of equidistant values in the range of the latency values
    x = np.linspace(min(latency_list), max(latency_list), (max(latency_list) - min(latency_list)) * 10)
    y = ag(x)  # evaluate the latency values for each x value
    return x, y

# draw_histogram(S)
x, y = get_pdf(S)
plt.rcParams.update({'font.size': 13.45})
plt.xlabel('Step (S)')
plt.legend()
plt.plot(x, y, 'r')
plt.show()
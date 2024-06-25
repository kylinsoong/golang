#!/usr/bin/python3

import sys
import re
import matplotlib.pyplot as plt

class AS3BenchMark:
    def __init__(self, ops, runTime, totalTime):
        self.ops = ops
        self.runTime = runTime
        self.totalTime = totalTime

    def __str__(self):
        return f"Operation: {self.ops}, Run Time: {self.runTime}s, Total Time: {self.totalTime}s"

class AS3BenchMarkResults(AS3BenchMark):
    def __init__(self, ops, runTime, totalTime, totalService):
        super().__init__(ops, runTime, totalTime)
        self.totalService = totalService

def extract_blocks(log_data, block_type):
    pattern = rf"{block_type} \d+\n(?:\d{{4}}.*\n)+"
    matches = re.findall(pattern, log_data, re.MULTILINE)
    return matches

def trip_prefix(line, prefix):
    if len(line) > 0 and prefix is not None and prefix in line:
        return line.strip().lstrip(prefix).strip()
    else:
        return line.strip()

def plot_benchmarks(title, benchmarks):
    operations = [benchmark.totalService for benchmark in benchmarks]
    runtime = [benchmark.runTime for benchmark in benchmarks]
    totaltime = [benchmark.totalTime for benchmark in benchmarks]

    # Plotting
    plt.plot(operations, runtime, marker='o', label='Runtime')
    plt.plot(operations, totaltime, marker='s', label='Total Time')

    # Adding labels and title
    plt.title(title)
    plt.xlabel('Service Number')
    plt.ylabel('Time (s)')
    plt.xticks(rotation=45)
    plt.legend()

    # Display the plot
    plt.tight_layout()
    plt.show()

def load_benchmark_logging(fileconfig):
    with open(fileconfig, 'r') as fo:
        data_all = fo.read()
    fo.close()
    return data_all

if not sys.argv[1:]:
    print("Usage: one-tenant.py [file]")
    sys.exit()

fileconfig = sys.argv[1]

benchmark_logging = load_benchmark_logging(fileconfig)

add_blocks = extract_blocks(benchmark_logging, "ADD")
del_blocks = extract_blocks(benchmark_logging, "DEL")
benchmarks = []
add_benchmarks = []
del_benchmarks = []
benchmarks_by_operation = {}
blocks = add_blocks + del_blocks

print(type(add_blocks), type(del_blocks))

for block in blocks:
    lines = block.strip().split("\n")
    ops = lines[0]
    parts = lines[3].split(',')
    runtime = trip_prefix(parts[-1], "runTime:")
    parts = lines[4].split(':')
    totaltime = parts[-1].strip().rstrip('s')
    runtime = int(runtime) / 1000
    totaltime = "{:.3f}".format(float(totaltime))
    benchmark = AS3BenchMark(ops, runtime, float(totaltime))
    benchmarks.append(benchmark)

for benchmark in benchmarks:
    if benchmark.ops not in benchmarks_by_operation:
        benchmarks_by_operation[benchmark.ops] = []
    benchmarks_by_operation[benchmark.ops].append(benchmark)

for operation, benchmarks_list in benchmarks_by_operation.items():
    benchmarks_list.sort(key=lambda x: x.totalTime)
    benchmarks_list = benchmarks_list[2:-2]
    total_benchmarks = len(benchmarks_list)
    avg_runtime = sum(benchmark.runTime for benchmark in benchmarks_list) / total_benchmarks
    avg_totaltime = sum(benchmark.totalTime for benchmark in benchmarks_list) / total_benchmarks
    
    parts = operation.split()
    results = AS3BenchMarkResults(operation, "{:.3f}".format(avg_runtime), "{:.3f}".format(avg_totaltime), int(parts[1]))
    if parts[0] == "ADD":
        add_benchmarks.append(results)
    else:
        del_benchmarks.append(results)
   
add_benchmarks.sort(key=lambda x: x.totalService)   
del_benchmarks.sort(key=lambda x: x.totalService)   
add_benchmarks = add_benchmarks[1:]
del_benchmarks = del_benchmarks[1:]

for benchmark in add_benchmarks:
    print(f"{benchmark.ops}, runtime: {benchmark.runTime}s, totaltime: {benchmark.totalTime}s")

for benchmark in del_benchmarks:
    print(f"{benchmark.ops}, runtime: {benchmark.runTime}s, totaltime: {benchmark.totalTime}s")

plot_benchmarks("AS3 ADD Operation", add_benchmarks)
plot_benchmarks("AS3 DELETE Operation", del_benchmarks)

K = int.from_bytes(b"flag{ch1n3s3-r3mAind3r-D-l0g}")

primes = [191, 1621, 61, 2447, 991, 1297, 47, 1049, 347, 283, 2617, 1429, 167, 307, 431, 683, 1627, 17, 827, 97, 523, 151, 37, 2269, 1733, 3, 19, 439]
d_logs = []

for q in primes:
    print("computing discrete log for prime: " + str(q))
    d_log = K % q
    d_logs.append(d_log)
    print("discrete log found: " + str(d_log))

print(d_logs)
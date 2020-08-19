# https://en.bitcoin.it/wiki/Secp256k1
_p = 115792089237316195423570985008687907853269984665640564039457584007908834671663
_n = 115792089237316195423570985008687907852837564279074904382605163141518161494337
_a = 0
_b = 7
_gx = int("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)
_gy = int("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)
_g = (_gx, _gy)


# https://en.wikipedia.org/wiki/Extended_Euclidean_algorithm
def inv_mod(a, n=_p):
    lm, hm = 1, 0
    low, high = a % n, n
    while low > 1:
        ratio = high // low
        nm, new = hm - lm * ratio, high - low * ratio
        lm, low, hm, high = nm, new, lm, low
    return lm % n


# https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Point_addition
def ecc_add(a, b):
    l = ((b[1] - a[1]) * inv_mod(b[0] - a[0])) % _p
    x = (l * l - a[0] - b[0]) % _p
    y = (l * (a[0] - x) - a[1]) % _p
    return (x, y)


# https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Point_doubling
def ecc_double(a):
    l = ((3 * a[0] * a[0] + _a) * inv_mod((2 * a[1]))) % _p
    x = (l * l - 2 * a[0]) % _p
    y = (l * (a[0] - x) - a[1]) % _p
    return (x, y)


# https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication
def ecc_mul(point, scalar):
    if scalar == 0 or scalar >= _p:
        raise ValueError("INVALID_SCALAR_OR_PRIVATEKEY")
    scalar_bin = str(bin(scalar))[2:]
    q = point
    for i in range(1, len(scalar_bin)):
        q = ecc_double(q)
        if scalar_bin[i] == "1":
            q = ecc_add(q, point)
    return q


# https://rosettacode.org/wiki/Cipolla%27s_algorithm#Python
def to_base(n, b):
    if n < 2:
        return [n]
    temp = n
    ans = []
    while temp != 0:
        ans = [temp % b] + ans
        temp //= b
    return ans


# https://rosettacode.org/wiki/Cipolla%27s_algorithm#Python
def ecc_sqrt(n, p):
    n %= p
    if n == 0 or n == 1:
        return (n, -n % p)
    phi = p - 1
    if pow(n, phi // 2, p) != 1:
        return ()
    if p % 4 == 3:
        ans = pow(n, (p + 1) // 4, p)
        return (ans, -ans % p)
    aa = 0
    for i in range(1, p):
        temp = pow((i * i - n) % p, phi // 2, p)
        if temp == phi:
            aa = i
            break
    exponent = to_base((p + 1) // 2, 2)

    def cipolla_mult(ab, cd, w, p):
        a, b = ab
        c, d = cd
        return ((a * c + b * d * w) % p, (a * d + b * c) % p)

    x1 = (aa, 1)
    x2 = cipolla_mult(x1, x1, aa * aa - n, p)
    for i in range(1, len(exponent)):
        if exponent[i] == 0:
            x2 = cipolla_mult(x2, x1, aa * aa - n, p)
            x1 = cipolla_mult(x1, x1, aa * aa - n, p)
        else:
            x1 = cipolla_mult(x1, x2, aa * aa - n, p)
            x2 = cipolla_mult(x2, x2, aa * aa - n, p)

    return (x1[0], -x1[0] % p)


# https://en.wikipedia.org/wiki/Elliptic_Curve_Digital_Signature_Algorithm
def ecrecover(_e: bytes, _r: bytes, _s: bytes, v):
    e = int.from_bytes(_e, "big")
    r = int.from_bytes(_r, "big")
    s = int.from_bytes(_s, "big")

    x = r % _n
    y1, y2 = ecc_sqrt(x * x * x + x * _a + _b, _p)
    if v == 27:
        y = y1 if y1 % 2 == 0 else y2
    elif v == 28:
        y = y1 if y1 % 2 == 1 else y2
    else:
        raise ValueError(f"ECRECOVER_ERROR: v must be 27 or 28 but got {v}")

    R = (x, y % _n)
    x_inv = inv_mod(x, _n)
    gxh = ecc_mul(_g, -e % _n)

    pub = ecc_mul(ecc_add(gxh, ecc_mul(R, s)), x_inv)

    return bytes.fromhex("%064x" % pub[0] + "%064x" % pub[1])

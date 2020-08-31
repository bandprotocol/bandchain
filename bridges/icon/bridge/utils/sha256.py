# https://gist.github.com/prokls/41e82472bd4968720d1482f81235e0ac

F32 = 0xFFFFFFFF

_k = [0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5,
      0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
      0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3,
      0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
      0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc,
      0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
      0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7,
      0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
      0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13,
      0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
      0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3,
      0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
      0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5,
      0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
      0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208,
      0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2]

_h = [0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a,
      0x510e527f, 0x9b05688c, 0x1f83d9ab, 0x5be0cd19]


def _pad(msglen):
    mdi = msglen & 0x3F
    length = (msglen << 3).to_bytes(8, byteorder='big')

    if mdi < 56:
        padlen = 55 - mdi
    else:
        padlen = 119 - mdi

    return b'\x80' + (b'\x00' * padlen) + length


def _rotr(x, y):
    return ((x >> y) | (x << (32 - y))) & F32


def _maj(x, y, z):
    return (x & y) ^ (x & z) ^ (y & z)


def _ch(x, y, z):
    return (x & y) ^ ((~x) & z)


def _compress(c, hh):
    k = _k[:]
    w = [0] * 64
    w[0:16] = tuple([int.from_bytes(c[i*4:i*4+4], byteorder='big')
                     for i in range(16)])

    for i in range(16, 64):
        s0 = _rotr(w[i-15], 7) ^ _rotr(w[i-15], 18) ^ (w[i-15] >> 3)
        s1 = _rotr(w[i-2], 17) ^ _rotr(w[i-2], 19) ^ (w[i-2] >> 10)
        w[i] = (w[i-16] + s0 + w[i-7] + s1) & F32

    a, b, c, d, e, f, g, h = hh

    for i in range(64):
        s0 = _rotr(a, 2) ^ _rotr(a, 13) ^ _rotr(a, 22)
        t2 = s0 + _maj(a, b, c)
        s1 = _rotr(e, 6) ^ _rotr(e, 11) ^ _rotr(e, 25)
        t1 = h + s1 + _ch(e, f, g) + k[i] + w[i]

        h = g
        g = f
        f = e
        e = (d + t1) & F32
        d = c
        c = b
        b = a
        a = (t1 + t2) & F32

    for i, (x, y) in enumerate(zip(hh, [a, b, c, d, e, f, g, h])):
        hh[i] = (x + y) & F32

    return hh


def update(counter, cache, m, h):
    if not m:
        return counter, cache, h

    counter += len(m)
    m = cache + m

    for i in range(0, len(m) // 64):
        h = _compress(m[64 * i:64 * (i + 1)], h)
    cache = m[-(len(m) % 64):]

    return counter, cache, h


def digest(_m):
    counter, cache, h = update(0, b'', _m, _h[:])
    counter, cache, h = update(counter, cache, _pad(counter), h)
    return b''.join([(i).to_bytes(4, byteorder='big') for i in h[:8]])

import setuptools

with open("README.md", "r") as fh:
    long_description = fh.read()

setuptools.setup(
    name="pyband",
    packages=["pyband"],
    version="0.1.0",
    license="MIT",
    description="Python library for BandChain",
    long_description=long_description,
    long_description_content_type="text/markdown",
    author="Band Protocol",
    author_email="dev@bandprotocol.com",
    url="https://github.com/bandprotocol/bandchain",
    keywords=["BAND", "BLOCKCHAIN", "ORACLE"],
    install_requires=["requests", "dacite", "bech32", "bip32", "ecdsa", "mnemonic"],
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.6",
        "Programming Language :: Python :: 3.7",
        "Programming Language :: Python :: 3.8",
    ],
)

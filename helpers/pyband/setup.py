from distutils.core import setup

setup(
    name="pyband",
    packages=["pyband"],
    version="0.0.3",
    license="MIT",
    description="Python library for BandChain",
    author="Band Protocol",
    author_email="dev@bandprotocol.com",
    url="https://github.com/bandprotocol/bandchain",
    keywords=["BAND", "BLOCKCHAIN", "ORACLE"],
    install_requires=["requests", "dacite"],
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

// SPDX-License-Identifier: Apache-2.0

pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;

import {Ownable} from "openzeppelin-solidity/contracts/access/Ownable.sol";
import {SafeMath} from "openzeppelin-solidity/contracts/math/SafeMath.sol";
import {Strings} from "./libraries/Strings.sol";
import {ICacheBridge} from "./interfaces/ICacheBridge.sol";
import {IBandDataset} from "./interfaces/IBandDataset.sol";
import {ParamsDecoder} from "./libraries/BandStandardParamsDecoder.sol";
import {ResultDecoder} from "./libraries/BandStandardResultDecoder.sol";

contract BandDataset is IBandDataset, Ownable {
    using Strings for *;
    using SafeMath for uint256;

    using ResultDecoder for bytes;
    using ParamsDecoder for bytes;

    ICacheBridge bridge;

    mapping(string => SymbolData) dataFromSymbol;

    bytes[] calldataArray = new bytes[](6);

    constructor(ICacheBridge _bridge) public {
        bridge = _bridge;

        // ["BTC","ETH","USDT","XRP","LINK","DOT","BCH","LTC","ADA","BSV","CRO","BNB","EOS","XTZ","TRX","XLM","ATOM","XMR","OKB","USDC","NEO","XEM","LEO","HT","VET"]
        calldataArray[0] = hex"000000190000000342544300000003455448000000045553445400000003585250000000044c494e4b00000003444f5400000003424348000000034c544300000003414441000000034253560000000343524f00000003424e4200000003454f530000000358545a0000000354525800000003584c4d0000000441544f4d00000003584d52000000034f4b420000000455534443000000034e454f0000000358454d000000034c454f00000002485400000003564554000000003b9aca00";
        // ["YFI","MIOTA","LEND","SNX","DASH","COMP","ZEC","ETC","OMG","MKR","ONT","NXM","AMPL","BAT","THETA","DAI","REN","ZRX","ALGO","FTT","DOGE","KSM","WAVES","EWT","DGB"]
        calldataArray[1] = hex"0000001900000003594649000000054d494f5441000000044c454e4400000003534e58000000044441534800000004434f4d50000000035a454300000003455443000000034f4d47000000034d4b52000000034f4e54000000034e584d00000004414d504c00000003424154000000055448455441000000034441490000000352454e000000035a525800000004414c474f0000000346545400000004444f4745000000034b534d0000000557415645530000000345575400000003444742000000003b9aca00";
        // ["KNC","ICX","TUSD","SUSHI","BTT","BAND","EGLD","ANT","NMR","PAX","LSK","LRC","HBAR","BAL","RUNE","YFII","LUNA","DCR","SC","STX","ENJ","BUSD","OCEAN","RSR","SXP"]
        calldataArray[2] = hex"00000019000000034b4e43000000034943580000000454555344000000055355534849000000034254540000000442414e440000000445474c4400000003414e54000000034e4d5200000003504158000000034c534b000000034c524300000004484241520000000342414c0000000452554e450000000459464949000000044c554e41000000034443520000000253430000000353545800000003454e4a0000000442555344000000054f4345414e0000000352535200000003535850000000003b9aca00";
        // ["BTG","BZRX","SRM","SNT","SOL","CKB","BNT","CRV","MANA","YFV","KAVA","MATIC","TRB","REP","FTM","TOMO","ONE","WNXM","PAXG","WAN","SUSD","RLC","OXT","RVN","FNX"]
        calldataArray[3] = hex"000000190000000342544700000004425a52580000000353524d00000003534e5400000003534f4c00000003434b4200000003424e5400000003435256000000044d414e4100000003594656000000044b415641000000054d4154494300000003545242000000035245500000000346544d00000004544f4d4f000000034f4e4500000004574e584d00000004504158470000000357414e000000045355534400000003524c43000000034f58540000000352564e00000003464e58000000003b9aca00";
        // ["EUR","GBP","CNY","SGD","RMB","KRW","JPY","INR","RUB","CHF","AUD","BRL","CAD","HKD","XAU","XAG"]
        calldataArray[4] = hex"00000010000000034555520000000347425000000003434e590000000353474400000003524d42000000034b5257000000034a505900000003494e520000000352554200000003434846000000034155440000000342524c0000000343414400000003484b440000000358415500000003584147000000003b9aca00";
        // ["RENBTC","WBTC","DIA","BTM","IOTX","FET","JST","MCO","KMD","BTS","QKC","YAMV2","XZC","UOS","AKRO","HNT","HOT","KAI","OGN","WRX","KDA","ORN","FOR","AST","STORJ"]
        calldataArray[5] = hex"000000190000000652454e4254430000000457425443000000034449410000000342544d00000004494f545800000003464554000000034a5354000000034d434f000000034b4d440000000342545300000003514b430000000559414d563200000003585a4300000003554f5300000004414b524f00000003484e5400000003484f54000000034b4149000000034f474e00000003575258000000034b4441000000034f524e00000003464f52000000034153540000000553544f524a000000003b9aca00";

        // Tokens
        dataFromSymbol["BTC"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 0
        });
        dataFromSymbol["ETH"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 1
        });
        dataFromSymbol["USDT"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 2
        });
        dataFromSymbol["XRP"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 3
        });
        dataFromSymbol["LINK"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 4
        });
        dataFromSymbol["DOT"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 5
        });
        dataFromSymbol["BCH"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 6
        });
        dataFromSymbol["LTC"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 7
        });
        dataFromSymbol["ADA"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 8
        });
        dataFromSymbol["BSV"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 9
        });
        dataFromSymbol["CRO"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 10
        });
        dataFromSymbol["BNB"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 11
        });
        dataFromSymbol["EOS"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 12
        });
        dataFromSymbol["XTZ"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 13
        });
        dataFromSymbol["TRX"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 14
        });
        dataFromSymbol["XLM"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 15
        });
        dataFromSymbol["ATOM"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 16
        });
        dataFromSymbol["XMR"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 17
        });
        dataFromSymbol["OKB"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 18
        });
        dataFromSymbol["USDC"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 19
        });
        dataFromSymbol["NEO"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 20
        });
        dataFromSymbol["XEM"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 21
        });
        dataFromSymbol["LEO"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 22
        });
        dataFromSymbol["HT"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 23
        });
        dataFromSymbol["VET"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 0,
            rateID: 24
        });

        dataFromSymbol["YFI"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 0
        });
        dataFromSymbol["MIOTA"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 1
        });
        dataFromSymbol["LEND"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 2
        });
        dataFromSymbol["SNX"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 3
        });
        dataFromSymbol["DASH"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 4
        });
        dataFromSymbol["COMP"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 5
        });
        dataFromSymbol["ZEC"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 6
        });
        dataFromSymbol["ETC"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 7
        });
        dataFromSymbol["OMG"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 8
        });
        dataFromSymbol["MKR"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 9
        });
        dataFromSymbol["ONT"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 10
        });
        dataFromSymbol["NXM"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 11
        });
        dataFromSymbol["AMPL"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 12
        });
        dataFromSymbol["BAT"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 13
        });
        dataFromSymbol["THETA"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 14
        });
        dataFromSymbol["DAI"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 15
        });
        dataFromSymbol["REN"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 16
        });
        dataFromSymbol["ZRX"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 17
        });
        dataFromSymbol["ALGO"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 18
        });
        dataFromSymbol["FTT"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 19
        });
        dataFromSymbol["DOGE"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 20
        });
        dataFromSymbol["KSM"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 21
        });
        dataFromSymbol["WAVES"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 22
        });
        dataFromSymbol["EWT"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 23
        });
        dataFromSymbol["DGB"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 1,
            rateID: 24
        });

        dataFromSymbol["KNC"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 0
        });
        dataFromSymbol["ICX"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 1
        });
        dataFromSymbol["TUSD"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 2
        });
        dataFromSymbol["SUSHI"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 3
        });
        dataFromSymbol["BTT"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 4
        });
        dataFromSymbol["BAND"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 5
        });
        dataFromSymbol["ERD"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 6
        });
        dataFromSymbol["ANT"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 7
        });
        dataFromSymbol["NMR"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 8
        });
        dataFromSymbol["PAX"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 9
        });
        dataFromSymbol["LSK"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 10
        });
        dataFromSymbol["LRC"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 11
        });
        dataFromSymbol["HBAR"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 12
        });
        dataFromSymbol["BAL"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 13
        });
        dataFromSymbol["RUNE"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 14
        });
        dataFromSymbol["YFII"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 15
        });
        dataFromSymbol["LUNA"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 16
        });
        dataFromSymbol["DCR"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 17
        });
        dataFromSymbol["SC"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 18
        });
        dataFromSymbol["STX"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 19
        });
        dataFromSymbol["ENJ"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 20
        });
        dataFromSymbol["BUSD"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 21
        });
        dataFromSymbol["OCEAN"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 22
        });
        dataFromSymbol["RSR"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 23
        });
        dataFromSymbol["SXP"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 2,
            rateID: 24
        });

        dataFromSymbol["BTG"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 0
        });
        dataFromSymbol["BZRX"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 1
        });
        dataFromSymbol["SRM"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 2
        });
        dataFromSymbol["SNT"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 3
        });
        dataFromSymbol["SOL"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 4
        });
        dataFromSymbol["CKB"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 5
        });
        dataFromSymbol["BNT"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 6
        });
        dataFromSymbol["CRV"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 7
        });
        dataFromSymbol["MANA"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 8
        });
        dataFromSymbol["YFV"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 9
        });
        dataFromSymbol["KAVA"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 10
        });
        dataFromSymbol["MATIC"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 11
        });
        dataFromSymbol["TRB"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 12
        });
        dataFromSymbol["REP"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 13
        });
        dataFromSymbol["FTM"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 14
        });
        dataFromSymbol["TOMO"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 15
        });
        dataFromSymbol["ONE"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 16
        });
        dataFromSymbol["WNXM"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 17
        });
        dataFromSymbol["PAXG"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 18
        });
        dataFromSymbol["WAN"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 19
        });
        dataFromSymbol["SUSD"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 20
        });
        dataFromSymbol["RLC"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 21
        });
        dataFromSymbol["OXT"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 22
        });
        dataFromSymbol["RVN"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 23
        });
        dataFromSymbol["FNX"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 3,
            rateID: 24
        });

        dataFromSymbol["EUR"] = SymbolData({
            oracleScriptID: 9,
            calldataID: 4,
            rateID: 0
        });
        dataFromSymbol["GBP"] = SymbolData({
            oracleScriptID: 9,
            calldataID: 4,
            rateID: 1
        });
        dataFromSymbol["CNY"] = SymbolData({
            oracleScriptID: 9,
            calldataID: 4,
            rateID: 2
        });
        dataFromSymbol["SGD"] = SymbolData({
            oracleScriptID: 9,
            calldataID: 4,
            rateID: 3
        });
        dataFromSymbol["RMB"] = SymbolData({
            oracleScriptID: 9,
            calldataID: 4,
            rateID: 4
        });
        dataFromSymbol["KRW"] = SymbolData({
            oracleScriptID: 9,
            calldataID: 4,
            rateID: 5
        });
        dataFromSymbol["JPY"] = SymbolData({
            oracleScriptID: 9,
            calldataID: 4,
            rateID: 6
        });
        dataFromSymbol["INR"] = SymbolData({
            oracleScriptID: 9,
            calldataID: 4,
            rateID: 7
        });
        dataFromSymbol["RUB"] = SymbolData({
            oracleScriptID: 9,
            calldataID: 4,
            rateID: 8
        });
        dataFromSymbol["CHF"] = SymbolData({
            oracleScriptID: 9,
            calldataID: 4,
            rateID: 9
        });
        dataFromSymbol["AUD"] = SymbolData({
            oracleScriptID: 9,
            calldataID: 4,
            rateID: 10
        });
        dataFromSymbol["BRL"] = SymbolData({
            oracleScriptID: 9,
            calldataID: 4,
            rateID: 11
        });
        dataFromSymbol["CAD"] = SymbolData({
            oracleScriptID: 9,
            calldataID: 4,
            rateID: 12
        });
        dataFromSymbol["HKD"] = SymbolData({
            oracleScriptID: 9,
            calldataID: 4,
            rateID: 13
        });
        dataFromSymbol["XAU"] = SymbolData({
            oracleScriptID: 9,
            calldataID: 4,
            rateID: 14
        });
        dataFromSymbol["XAG"] = SymbolData({
            oracleScriptID: 9,
            calldataID: 4,
            rateID: 15
        });

        dataFromSymbol["RENBTC"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 0
        });
        dataFromSymbol["WBTC"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 1
        });
        dataFromSymbol["DIA"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 2
        });
        dataFromSymbol["BTM"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 3
        });
        dataFromSymbol["IOTX"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 4
        });
        dataFromSymbol["FET"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 5
        });
        dataFromSymbol["JST"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 6
        });
        dataFromSymbol["MCO"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 7
        });
        dataFromSymbol["KMD"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 8
        });
        dataFromSymbol["BTS"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 9
        });
        dataFromSymbol["QKC"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 10
        });
        dataFromSymbol["YAMV2"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 11
        });
        dataFromSymbol["XZC"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 12
        });
        dataFromSymbol["UOS"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 13
        });
        dataFromSymbol["AKRO"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 14
        });
        dataFromSymbol["HNT"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 15
        });
        dataFromSymbol["HOT"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 16
        });
        dataFromSymbol["KAI"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 17
        });
        dataFromSymbol["OGN"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 18
        });
        dataFromSymbol["WRX"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 19
        });
        dataFromSymbol["KDA"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 20
        });
        dataFromSymbol["ORN"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 21
        });
        dataFromSymbol["FOR"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 22
        });
        dataFromSymbol["AST"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 23
        });
        dataFromSymbol["STORJ"] = SymbolData({
            oracleScriptID: 8,
            calldataID: 5,
            rateID: 24
        });
    }

    function setBridge(ICacheBridge _bridge) external onlyOwner {
        bridge = _bridge;
    }

    function getReferenceData(string[] memory pairs)
        external
        override
        view
        returns (ReferenceData[] memory)
    {
        ICacheBridge.RequestPacket memory req;
        req.clientId = "bandteam";
        req.askCount = 4;
        req.minCount = 3;
        ReferenceData[] memory result = new ReferenceData[](pairs.length);

        ICacheBridge.ResponsePacket memory latestResponse;
        ResultDecoder.Result memory decodedResult;

        for (uint256 i = 0; i < pairs.length; i++) {
            DataUpdate memory updateTimes;

            Strings.slice memory s = pairs[i].toSlice();
            Strings.slice memory delim = "/".toSlice();

            // Base Symbol
            uint64 basePrice;
            string memory base = s.split(delim).toString();
            if (
                keccak256(abi.encodePacked(base)) ==
                keccak256(abi.encodePacked("USD"))
            ) {
                updateTimes.base = 0;
                basePrice = 1e9;
            } else {
                req.oracleScriptId = dataFromSymbol[base].oracleScriptID;
                req.params = calldataArray[dataFromSymbol[base].calldataID];

                latestResponse = bridge.getLatestResponse(req);
                decodedResult = latestResponse.result.decodeResult();
                updateTimes.base = latestResponse.resolveTime;

                basePrice = decodedResult.rates[dataFromSymbol[base].rateID];
            }

            // Quote Symbol
            uint64 quotePrice;
            string memory quote = s.split(delim).toString();
            if (
                keccak256(abi.encodePacked(quote)) ==
                keccak256(abi.encodePacked("USD"))
            ) {
                updateTimes.quote = 0;
                quotePrice = 1e9;
            } else {
                req.oracleScriptId = dataFromSymbol[quote].oracleScriptID;
                req.params = calldataArray[dataFromSymbol[quote].calldataID];

                latestResponse = bridge.getLatestResponse(req);
                decodedResult = latestResponse.result.decodeResult();
                updateTimes.quote = latestResponse.resolveTime;

                quotePrice = decodedResult.rates[dataFromSymbol[quote].rateID];
            }

            result[i].rate = (uint256(basePrice) * 1e18) / uint256(quotePrice);
            result[i].lastUpdated = updateTimes;
        }
        return result;
    }
}

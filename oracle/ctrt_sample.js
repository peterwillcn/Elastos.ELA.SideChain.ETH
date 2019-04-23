"use strict";

module.exports = {
    abi: [{
        "constant": true,
        "inputs": [],
        "name": "getSubContractAddress",
        "outputs": [{"name": "addr", "type": "address"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    }, {
        "constant": true,
        "inputs": [{"name": "", "type": "bytes32"}],
        "name": "txProcessed",
        "outputs": [{"name": "", "type": "bool"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    }, {
        "constant": false,
        "inputs": [{"name": "_txhash", "type": "bytes32"}],
        "name": "sendPayload",
        "outputs": [],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    }, {
        "constant": false,
        "inputs": [{"name": "_addrs", "type": "string[]"}, {"name": "_amounts", "type": "uint256[]"}, {
            "name": "_fee",
            "type": "uint256"
        }],
        "name": "receivePayload",
        "outputs": [],
        "payable": true,
        "stateMutability": "payable",
        "type": "function"
    }, {
        "constant": true,
        "inputs": [{"name": "", "type": "bytes32"}, {"name": "", "type": "uint256"}],
        "name": "amounts",
        "outputs": [{"name": "", "type": "uint256"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    }, {
        "constant": true,
        "inputs": [{"name": "", "type": "bytes32"}, {"name": "", "type": "uint256"}],
        "name": "indexes",
        "outputs": [{"name": "", "type": "uint256"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    }, {
        "constant": true,
        "inputs": [{"name": "", "type": "bytes32"}],
        "name": "num",
        "outputs": [{"name": "", "type": "uint256"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    }, {
        "constant": true,
        "inputs": [{"name": "", "type": "bytes32"}, {"name": "", "type": "uint256"}],
        "name": "addrs",
        "outputs": [{"name": "", "type": "address"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    }, {"payable": true, "stateMutability": "payable", "type": "fallback"}, {
        "anonymous": false,
        "inputs": [{"indexed": true, "name": "_txhash", "type": "bytes32"}, {
            "indexed": true,
            "name": "_addr",
            "type": "address"
        }, {"indexed": false, "name": "_amount", "type": "uint256"}, {
            "indexed": true,
            "name": "_arbiter",
            "type": "address"
        }],
        "name": "PayloadSent",
        "type": "event"
    }, {
        "anonymous": false,
        "inputs": [{"indexed": true, "name": "_txhash", "type": "bytes32"}, {
            "indexed": true,
            "name": "_arbiter",
            "type": "address"
        }],
        "name": "TxProcessed",
        "type": "event"
    }, {
        "anonymous": false,
        "inputs": [{"indexed": false, "name": "_addr", "type": "string"}, {
            "indexed": false,
            "name": "_amount",
            "type": "uint256"
        }, {"indexed": false, "name": "_crosschainamount", "type": "uint256"}, {
            "indexed": true,
            "name": "_sender",
            "type": "address"
        }],
        "name": "PayloadReceived",
        "type": "event"
    }, {
        "anonymous": false,
        "inputs": [{"indexed": true, "name": "_sender", "type": "address"}, {
            "indexed": false,
            "name": "_amount",
            "type": "uint256"
        }],
        "name": "EtherDeposited",
        "type": "event"
    }],
    address: "0x491bC043672B9286fA02FA7e0d6A3E5A0384A31A"
}

// SPDX-License-Identifier: MIT
pragma solidity 0.8.20;

contract FlareDataConnectorHub {
    event AttestationRequest(bytes data, uint256 fee);

    uint256 public constant minimalFee = 1 wei;

    function requestAttestation(bytes calldata _data) external payable {
        require(msg.value >= minimalFee, "fee to low");

        emit AttestationRequest(_data, msg.value);
    }

    function getBaseFee(bytes calldata _data) external view returns (uint256) {
        return minimalFee;
    }
}

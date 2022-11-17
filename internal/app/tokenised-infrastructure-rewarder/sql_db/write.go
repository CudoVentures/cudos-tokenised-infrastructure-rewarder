package sql_db

import (
	"context"
	"github.com/CudoVentures/tokenised-infrastructure-rewarder/internal/app/tokenised-infrastructure-rewarder/types"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/jmoiron/sqlx"
)

func saveDestinationAddressesWithAmountHistory(ctx context.Context, tx *sqlx.Tx, address string, amountInfo types.AmountInfo, txHash string, farmId string) error {
	now := time.Now()
	if !amountInfo.ThresholdReached {
		txHash = "" // the funds were not sent but accumulated, we keep this record as statistic that they were spread but with empty tx hash
	}
	_, err := tx.ExecContext(ctx, insertDestinationAddressesWithAmountHistory, address, amountInfo.Amount, amountInfo.ThresholdReached, txHash, farmId, now.Unix(), now.UTC(), now.UTC())
	return err

}

// add to this table
func saveNFTInformationHistory(ctx context.Context, tx *sqlx.Tx, collectionDenomId, tokenId string,
	payoutPeriodStart, payoutPeriodEnd int64, reward btcutil.Amount, txHash string,
	maintenanceFee, CudoPartOfMaintenanceFee btcutil.Amount) error {
	now := time.Now()
	_, err := tx.ExecContext(ctx, insertNFTInformationHistory, collectionDenomId, tokenId, payoutPeriodStart,
		payoutPeriodEnd, reward, txHash, maintenanceFee, CudoPartOfMaintenanceFee, now.UTC(), now.UTC())
	return err
}

// add to this table
func saveNFTOwnersForPeriodHistory(ctx context.Context, tx *sqlx.Tx, collectionDenomId string, tokenId string, timedOwnedFrom int64,
	timedOwnedTo int64, totalTimeOwned int64, percentOfTimeOwned float64, owner string, payoutAddress string, reward btcutil.Amount) error {
	now := time.Now()
	_, err := tx.ExecContext(ctx, insertNFTOnwersForPeriodHistory, collectionDenomId, tokenId,
		timedOwnedFrom, timedOwnedTo, totalTimeOwned, percentOfTimeOwned, owner, payoutAddress, reward, now.UTC(), now.UTC())
	return err
}

func (sdb *SqlDB) SaveRBFTransactionHistory(ctx context.Context, tx *sqlx.Tx, oldTxHash string, newTxHash string, farm_id string) error {
	now := time.Now()
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, insertRBFTransactionHistory, oldTxHash, newTxHash,
			farm_id, now.UTC(), now.UTC())
	} else {
		_, err = sdb.db.ExecContext(ctx, insertRBFTransactionHistory, oldTxHash, newTxHash,
			farm_id, now.UTC(), now.UTC())
	}

	return err
}

func (sdb *SqlDB) SaveTxHashWithStatus(ctx context.Context, tx *sqlx.Tx, txHash string, txStatus string, farmSubAccountName string, retryCount int) error {
	now := time.Now()
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, insertTxHashWithStatus, txHash, txStatus, farmSubAccountName, retryCount, now.Unix(), now.UTC(), now.UTC())
	} else {
		_, err = sdb.db.ExecContext(ctx, insertTxHashWithStatus, txHash, txStatus, farmSubAccountName, retryCount, now.Unix(), now.UTC(), now.UTC())
	}
	return err
}

func (sdb *SqlDB) UpdateTransactionsStatus(ctx context.Context, tx *sqlx.Tx, txHashes []string, txStatus string) error {
	qry, args, err := sqlx.In(updateTxHashesWithStatusQuery, txStatus, txHashes)
	if err != nil {
		return err
	}

	if tx != nil {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = sdb.db.ExecContext(ctx, qry, args...)
	}
	if err != nil {
		return err
	}
	return nil
}

func (sdb *SqlDB) updateCurrentAcummulatedAmountForAddress(ctx context.Context, tx *sqlx.Tx, address string, farmId int, amount int64) error {
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, updateThresholdAmounts, amount, address, farmId)
	} else {
		_, err = sdb.db.ExecContext(ctx, updateThresholdAmounts, amount, address, farmId)
	}
	return err
}

func (sdb *SqlDB) markUTXOsAsProcessed(ctx context.Context, tx *sqlx.Tx, tx_hashes []string) interface{} {
	var UTXOMaps []map[string]interface{}
	for _, hash := range tx_hashes {
		m := map[string]interface{}{
			"tx_hash":   hash,
			"processed": true,
			"createdAt": time.Now().UTC(),
			"updatedAt": time.Now().UTC(),
		}
		UTXOMaps = append(UTXOMaps, m)
	}

	var err error
	if tx != nil {
		_, err = tx.NamedExec(insertUTXOWithStatus, UTXOMaps)
	} else {
		_, err = tx.NamedExecContext(ctx, insertUTXOWithStatus, UTXOMaps)
	}

	if err != nil {
		return err
	}
	return nil
}

func (sdb *SqlDB) SetInitialAccumulatedAmountForAddress(ctx context.Context, tx *sqlx.Tx, address string, farmId int, amount int) error {

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, insertInitialThresholdAmount, address, farmId, amount, time.Now().UTC(), time.Now().UTC())
	} else {
		_, err = sdb.db.ExecContext(ctx, insertInitialThresholdAmount, address, farmId, amount, time.Now().UTC(), time.Now().UTC())
	}
	return err

}

const (
	insertUTXOWithStatus = `INSERT INTO utxo_transactions (tx_hash, processed, "createdAt", "updatedAt")
	   VALUES (:tx_hash, :processed, :createdAt, :updatedAt)`

	insertTxHashWithStatus = `INSERT INTO statistics_tx_hash_status
	(tx_hash, status, time_sent, farm_sub_account_name, retry_count, "createdAt", "updatedAt") VALUES ($1, $2, $3, $4, $5, $6, $7)`
	insertRBFTransactionHistory = `INSERT INTO rbf_transaction_history
	(old_tx_hash, new_tx_hash, farm_sub_account_name, createdAt, updatedAt) VALUES ($1, $2, $3, $4, $5)`

	insertDestinationAddressesWithAmountHistory = `INSERT INTO statistics_destination_addresses_with_amount
		(address, amount, tx_hash, farm_id, payout_time, threshold_reached, "createdAt", "updatedAt") VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	insertNFTInformationHistory = `INSERT INTO statistics_nft_payout_history (denom_id, token_id, payout_period_start,
		payout_period_end, reward, tx_hash, maintenance_fee, cudo_part_of_maintenance_fee, "createdAt", "updatedAt")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	insertNFTOnwersForPeriodHistory = `INSERT INTO statistics_nft_owners_payout_history (denom_id, token_id, time_owned_from, time_owned_to,
		total_time_owned, percent_of_time_owned ,owner, payout_address, reward, "createdAt", "updatedAt")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	updateTxHashesWithStatusQuery = `UPDATE statistics_tx_hash_status SET status=? where tx_hash IN (?)`

	updateThresholdAmounts = `UPDATE threshold_amounts SET amount=$1 where btc_address=$2 and farm_id=$3`

	insertInitialThresholdAmount = `INSERT INTO threshold_amounts
	(btc_address, farm_id, amount, "createdAt", "updatedAt") VALUES ($1, $2, $3, $4, $5)`
)

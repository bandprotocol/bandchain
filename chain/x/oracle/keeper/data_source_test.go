package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func TestHasDataSource(t *testing.T) {
	_, ctx, k := createTestInput()
	// We should not have a data source ID 42 without setting it.
	require.False(t, k.HasDataSource(ctx, 42))
	// After we set it, we should be able to find it.
	k.SetDataSource(ctx, 42, types.NewDataSource(
		Owner.Address, BasicName, BasicDesc, Coins10uband, BasicExec,
	))
	require.True(t, k.HasDataSource(ctx, 42))
}

func TestSetterGetterDataSource(t *testing.T) {
	_, ctx, k := createTestInput()
	// Getting a non-existent data source should return error.
	_, err := k.GetDataSource(ctx, 42)
	require.Error(t, err)
	require.Panics(t, func() { _ = k.MustGetDataSource(ctx, 42) })
	// Creates some basic data sources.
	dataSource1 := types.NewDataSource(
		Alice.Address, "NAME1", "DESCRIPTION1", Coins10uband, []byte("executable1"),
	)
	dataSource2 := types.NewDataSource(
		Bob.Address, "NAME2", "DESCRIPTION2", Coins10uband, []byte("executable2"),
	)
	// Sets id 42 with data soure 1 and id 42 with data source 2.
	k.SetDataSource(ctx, 42, dataSource1)
	k.SetDataSource(ctx, 43, dataSource2)
	// Checks that Get and MustGet perform correctly.
	dataSource1Res, err := k.GetDataSource(ctx, 42)
	require.Nil(t, err)
	require.Equal(t, dataSource1, dataSource1Res)
	require.Equal(t, dataSource1, k.MustGetDataSource(ctx, 42))
	dataSource2Res, err := k.GetDataSource(ctx, 43)
	require.Nil(t, err)
	require.Equal(t, dataSource2, dataSource2Res)
	require.Equal(t, dataSource2, k.MustGetDataSource(ctx, 43))
	// Replaces id 42 with another data source.
	k.SetDataSource(ctx, 42, dataSource2)
	require.NotEqual(t, dataSource1, k.MustGetDataSource(ctx, 42))
	require.Equal(t, dataSource2, k.MustGetDataSource(ctx, 42))
}

func TestAddDataSourceEditDataSourceBasic(t *testing.T) {
	_, ctx, k := createTestInput()
	// Creates some basic data sources.
	dataSource1 := types.NewDataSource(
		Alice.Address, "NAME1", "DESCRIPTION1", Coins10uband, []byte("executable1"),
	)
	dataSource2 := types.NewDataSource(
		Bob.Address, "NAME2", "DESCRIPTION2", Coins10uband, []byte("executable2"),
	)
	// Adds a new data source to the store. We should be able to retreive it back.
	id, err := k.AddDataSource(ctx, types.NewDataSource(
		dataSource1.Owner, dataSource1.Name, dataSource1.Description,
		dataSource1.Fee, dataSource1.Executable,
	))
	require.Nil(t, err)
	require.Equal(t, dataSource1, k.MustGetDataSource(ctx, id))
	require.NotEqual(t, dataSource2, k.MustGetDataSource(ctx, id))
	// Edits the data source. We should get the updated data source.
	err = k.EditDataSource(ctx, id, types.NewDataSource(
		dataSource2.Owner, dataSource2.Name, dataSource2.Description,
		dataSource2.Fee, dataSource2.Executable,
	))
	require.Nil(t, err)
	require.NotEqual(t, dataSource1, k.MustGetDataSource(ctx, id))
	require.Equal(t, dataSource2, k.MustGetDataSource(ctx, id))
}

func TestAddDataSourceDataSourceMustReturnCorrectID(t *testing.T) {
	_, ctx, k := createTestInput()
	// Initially we expect the data source count to be zero.
	count := k.GetDataSourceCount(ctx)
	require.Equal(t, count, int64(0))
	// Every new data source we add should return a new ID.
	id1, err := k.AddDataSource(ctx, types.NewDataSource(
		Owner.Address, BasicName, BasicDesc, Coins10uband, BasicExec,
	))
	require.Nil(t, err)
	require.Equal(t, id1, types.DID(1))
	// Adds another data source so now ID should be 2.
	id2, err := k.AddDataSource(ctx, types.NewDataSource(
		Owner.Address, BasicName, BasicDesc, Coins10uband, BasicExec,
	))
	require.Nil(t, err)
	require.Equal(t, id2, types.DID(2))
	// Finally we expect the data source to increase to 2 since we added two data sources.
	count = k.GetDataSourceCount(ctx)
	require.Equal(t, count, int64(2))
}

func TestEditDataSourceNonExistentDataSource(t *testing.T) {
	_, ctx, k := createTestInput()
	// Editing a non-existent data source should return error.
	err := k.EditDataSource(ctx, 42, types.NewDataSource(
		Owner.Address, BasicName, BasicDesc, Coins10uband, BasicExec,
	))
	require.Error(t, err)
}

func TestAddDataSourceTooLongName(t *testing.T) {
	_, ctx, k := createTestInput()
	// Sets max name length to 9. We should fail to add data source with name length 10.
	k.SetParam(ctx, types.KeyMaxNameLength, 9)
	_, err := k.AddDataSource(ctx, types.NewDataSource(
		Owner.Address, "0123456789", BasicDesc, Coins10uband, BasicExec,
	))
	require.Error(t, err)
	// Sets max name length to 10. We should now be able to add the data source.
	k.SetParam(ctx, types.KeyMaxNameLength, 10)
	_, err = k.AddDataSource(ctx, types.NewDataSource(
		Owner.Address, "0123456789", BasicDesc, Coins10uband, BasicExec,
	))
	require.Nil(t, err)
}

func TestEditDataSourceTooLongName(t *testing.T) {
	_, ctx, k := createTestInput()
	id, err := k.AddDataSource(ctx, types.NewDataSource(
		Owner.Address, BasicName, BasicDesc, Coins10uband, BasicExec,
	))
	require.Nil(t, err)
	dataSource := k.MustGetDataSource(ctx, id)
	// Sets max name length to 9. We should fail to edit data source with name length 10.
	k.SetParam(ctx, types.KeyMaxNameLength, 9)
	err = k.EditDataSource(ctx, id, types.NewDataSource(
		dataSource.Owner, "0123456789", dataSource.Description,
		dataSource.Fee, dataSource.Executable,
	))
	require.Error(t, err)
	// Sets max name length to 10. We should now be able to edit the data source.
	k.SetParam(ctx, types.KeyMaxNameLength, 10)
	err = k.EditDataSource(ctx, id, types.NewDataSource(
		dataSource.Owner, "0123456789", dataSource.Description,
		dataSource.Fee, dataSource.Executable,
	))
	require.Nil(t, err)
}

func TestAddDataSourceTooLongDescription(t *testing.T) {
	_, ctx, k := createTestInput()
	// Sets max desc length to 41. We should fail to add data source with desc length 42.
	k.SetParam(ctx, types.KeyMaxDescriptionLength, 41)
	_, err := k.AddDataSource(ctx, types.NewDataSource(
		Owner.Address, BasicName, "________THIS_STRING_HAS_SIZE_OF_42________",
		Coins10uband, BasicExec,
	))
	require.Error(t, err)
	// Sets max desc length to 42. We should now be able to add the data source.
	k.SetParam(ctx, types.KeyMaxDescriptionLength, 42)
	_, err = k.AddDataSource(ctx, types.NewDataSource(
		Owner.Address, BasicName, "________THIS_STRING_HAS_SIZE_OF_42________",
		Coins10uband, BasicExec,
	))
	require.Nil(t, err)
}

func TestEditDataSourceTooLongDescription(t *testing.T) {
	_, ctx, k := createTestInput()
	id, err := k.AddDataSource(ctx, types.NewDataSource(
		Owner.Address, BasicName, BasicDesc, Coins10uband, BasicExec,
	))
	require.Nil(t, err)
	dataSource := k.MustGetDataSource(ctx, id)
	// Sets max desc length to 41. We should fail to edit data source with name length 42.
	k.SetParam(ctx, types.KeyMaxDescriptionLength, 41)
	err = k.EditDataSource(ctx, id, types.NewDataSource(
		dataSource.Owner, dataSource.Name, "________THIS_STRING_HAS_SIZE_OF_42________",
		dataSource.Fee, dataSource.Executable,
	))
	require.Error(t, err)
	// Sets max desc length to 42. We should now be able to edit the data source.
	k.SetParam(ctx, types.KeyMaxDescriptionLength, 42)
	err = k.EditDataSource(ctx, id, types.NewDataSource(
		dataSource.Owner, dataSource.Name, "________THIS_STRING_HAS_SIZE_OF_42________",
		dataSource.Fee, dataSource.Executable,
	))
	require.Nil(t, err)
}

func TestAddDataSourceTooBigExecutable(t *testing.T) {
	_, ctx, k := createTestInput()
	// Sets max executable size to 40. We should fail to add data source with exec size 42.
	k.SetParam(ctx, types.KeyMaxExecutableSize, 40)
	_, err := k.AddDataSource(ctx, types.NewDataSource(
		Owner.Address, BasicName, BasicDesc, Coins10uband,
		[]byte("________THIS_STRING_HAS_SIZE_OF_42________"),
	))
	require.Error(t, err)
	// Sets max executable size to 50. We should now be able to add the data source.
	k.SetParam(ctx, types.KeyMaxExecutableSize, 50)
	_, err = k.AddDataSource(ctx, types.NewDataSource(
		Owner.Address, BasicName, BasicDesc, Coins10uband,
		[]byte("________THIS_STRING_HAS_SIZE_OF_42________"),
	))
	require.Nil(t, err)
}

func TestEditDataSourceTooBigExecutable(t *testing.T) {
	_, ctx, k := createTestInput()
	id, err := k.AddDataSource(ctx, types.NewDataSource(
		Owner.Address, BasicName, BasicDesc, Coins10uband, BasicExec,
	))
	require.Nil(t, err)
	dataSource := k.MustGetDataSource(ctx, id)
	// Sets max executable size to 40. We should fail to edit data source with exec size 42.
	k.SetParam(ctx, types.KeyMaxExecutableSize, 40)
	err = k.EditDataSource(ctx, id, types.NewDataSource(
		dataSource.Owner, dataSource.Name, dataSource.Description, dataSource.Fee,
		[]byte("________THIS_STRING_HAS_SIZE_OF_42________"),
	))
	require.Error(t, err)
	// Sets max executable size to 50. We should now be able to edit the data source.
	k.SetParam(ctx, types.KeyMaxExecutableSize, 50)
	err = k.EditDataSource(ctx, id, types.NewDataSource(
		dataSource.Owner, dataSource.Name, dataSource.Description, dataSource.Fee,
		[]byte("________THIS_STRING_HAS_SIZE_OF_42________"),
	))
	require.Nil(t, err)
}

func TestPayDataSourceFeeSuccess(t *testing.T) {
	app, ctx, k := createTestInput()
	id, err := k.AddDataSource(ctx, types.NewDataSource(
		Owner.Address, BasicName, BasicDesc, Coins10uband, BasicExec,
	))
	require.Nil(t, err)
	// If we pay from Alice, funds 10uband should be transferred from Alice to Owner
	ownerCoins := app.BankKeeper.GetAllBalances(ctx, Owner.Address)
	aliceCoins := app.BankKeeper.GetAllBalances(ctx, Alice.Address)
	err = k.PayDataSourceFee(ctx, id, Alice.Address)
	require.Nil(t, err)
	ownerNewCoins := app.BankKeeper.GetAllBalances(ctx, Owner.Address)
	aliceNewCoins := app.BankKeeper.GetAllBalances(ctx, Alice.Address)
	require.Equal(t, ownerCoins.Add(Coins10uband...), ownerNewCoins)
	require.Equal(t, aliceCoins.Sub(Coins10uband), aliceNewCoins)
}

func TestPayDataSourceTooExpensive(t *testing.T) {
	_, ctx, k := createTestInput()
	id, err := k.AddDataSource(ctx, types.NewDataSource(
		Owner.Address, BasicName, BasicDesc, Coins10uband.Add(Coins1000000uband...), BasicExec,
	))
	require.Nil(t, err)
	// Alice should not have enough coins to pay for this.
	err = k.PayDataSourceFee(ctx, id, Alice.Address)
	require.Error(t, err)
}

func TestPayDataSourceFeeFromOwner(t *testing.T) {
	app, ctx, k := createTestInput()
	id, err := k.AddDataSource(ctx, types.NewDataSource(
		Owner.Address, BasicName, BasicDesc, Coins10uband.Add(Coins1000000uband...), BasicExec,
	))
	require.Nil(t, err)
	// Even though it's extremely expensive, the owner should not have any problem.
	ownerCoins := app.BankKeeper.GetAllBalances(ctx, Owner.Address)
	err = k.PayDataSourceFee(ctx, id, Owner.Address)
	require.Nil(t, err)
	ownerNewCoins := app.BankKeeper.GetAllBalances(ctx, Owner.Address)
	// And his balance should not change at all.
	require.Equal(t, ownerCoins, ownerNewCoins)
}

func TestPayDataSourceFeeFree(t *testing.T) {
	app, ctx, k := createTestInput()
	id, err := k.AddDataSource(ctx, types.NewDataSource(
		Owner.Address, BasicName, BasicDesc, CoinsZero, BasicExec,
	))
	require.Nil(t, err)
	// If the fee is zero, no one should have their balances changed.
	ownerCoins := app.BankKeeper.GetAllBalances(ctx, Owner.Address)
	aliceCoins := app.BankKeeper.GetAllBalances(ctx, Alice.Address)
	err = k.PayDataSourceFee(ctx, id, Alice.Address)
	require.Nil(t, err)
	ownerNewCoins := app.BankKeeper.GetAllBalances(ctx, Owner.Address)
	aliceNewCoins := app.BankKeeper.GetAllBalances(ctx, Alice.Address)
	require.Equal(t, ownerCoins, ownerNewCoins)
	require.Equal(t, aliceCoins, aliceNewCoins)
}

func TestPayDataSourceFeeNotExistent(t *testing.T) {
	_, ctx, k := createTestInput()
	// Paying fee using a non-existent data source should fail.
	err := k.PayDataSourceFee(ctx, 42, Alice.Address)
	require.Error(t, err)
}

func TestGetAllDataSources(t *testing.T) {
	_, ctx, k := createTestInput()
	// Sets the data sources to the storage.
	dataSources := []types.DataSource{
		types.NewDataSource(
			Alice.Address, "NAME1", "DESCRIPTION1", Coins10uband, []byte("executable1"),
		),
		types.NewDataSource(
			Bob.Address, "NAME2", "DESCRIPTION2", Coins10uband, []byte("executable2"),
		),
	}
	k.SetDataSource(ctx, 1, dataSources[0])
	k.SetDataSource(ctx, 2, dataSources[1])
	// We should now be able to get all the existing data sources.
	require.Equal(t, dataSources, k.GetAllDataSources(ctx))
}

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
		Owner.Address, BasicName, BasicDesc, BasicFilename,
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
		Alice.Address, "NAME1", "DESCRIPTION1", "filename1",
	)
	dataSource2 := types.NewDataSource(
		Bob.Address, "NAME2", "DESCRIPTION2", "filename2",
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
		Alice.Address, "NAME1", "DESCRIPTION1", "FILENAME1",
	)
	dataSource2 := types.NewDataSource(
		Bob.Address, "NAME2", "DESCRIPTION2", "FILENAME2",
	)
	// Adds a new data source to the store. We should be able to retreive it back.
	id := k.AddDataSource(ctx, dataSource1)
	require.Equal(t, dataSource1, k.MustGetDataSource(ctx, id))
	require.NotEqual(t, dataSource2, k.MustGetDataSource(ctx, id))
	// Edits the data source. We should get the updated data source.
	k.MustEditDataSource(ctx, id, types.NewDataSource(
		dataSource2.Owner, dataSource2.Name, dataSource2.Description, dataSource2.Filename,
	))
	require.NotEqual(t, dataSource1, k.MustGetDataSource(ctx, id))
	require.Equal(t, dataSource2, k.MustGetDataSource(ctx, id))
}

func TestEditDataSourceDoNotModify(t *testing.T) {
	_, ctx, k := createTestInput()
	// Creates some basic data sources.
	dataSource1 := types.NewDataSource(
		Alice.Address, "NAME1", "DESCRIPTION1", "FILENAME1",
	)
	dataSource2 := types.NewDataSource(
		Bob.Address, types.DoNotModify, types.DoNotModify, "FILENAME2",
	)
	// Adds a new data source to the store. We should be able to retreive it back.
	id := k.AddDataSource(ctx, dataSource1)
	require.Equal(t, dataSource1, k.MustGetDataSource(ctx, id))
	require.NotEqual(t, dataSource2, k.MustGetDataSource(ctx, id))
	// Edits the data source. We should get the updated data source.
	k.MustEditDataSource(ctx, id, dataSource2)
	dataSourceRes := k.MustGetDataSource(ctx, id)
	require.NotEqual(t, dataSourceRes, dataSource1)
	require.NotEqual(t, dataSourceRes, dataSource2)
	require.Equal(t, dataSourceRes.Owner, dataSource2.Owner)
	require.Equal(t, dataSourceRes.Name, dataSource1.Name)
	require.Equal(t, dataSourceRes.Description, dataSource1.Description)
	require.Equal(t, dataSourceRes.Filename, dataSource2.Filename)
}

func TestAddDataSourceDataSourceMustReturnCorrectID(t *testing.T) {
	_, ctx, k := createTestInput()
	// Initially we expect the data source count to be zero.
	count := k.GetDataSourceCount(ctx)
	require.Equal(t, count, int64(0))
	// Every new data source we add should return a new ID.
	id1 := k.AddDataSource(ctx, types.NewDataSource(Owner.Address, BasicName, BasicDesc, BasicFilename))
	require.Equal(t, id1, types.DataSourceID(1))
	// Adds another data source so now ID should be 2.
	id2 := k.AddDataSource(ctx, types.NewDataSource(Owner.Address, BasicName, BasicDesc, BasicFilename))
	require.Equal(t, id2, types.DataSourceID(2))
	// Finally we expect the data source to increase to 2 since we added two data sources.
	count = k.GetDataSourceCount(ctx)
	require.Equal(t, count, int64(2))
}

// func TestEditDataSourceNonExistentDataSource(t *testing.T) {
// 	_, ctx, k := createTestInput()
// 	// Editing a non-existent data source should return error.
// 	err := k.EditDataSource(ctx, 42, types.NewDataSource(
// 		Owner.Address, BasicName, BasicDesc, BasicExec,
// 	))
// 	require.Error(t, err)
// }

// func TestGetAllDataSources(t *testing.T) {
// 	_, ctx, k := createTestInput()
// 	// Sets the data sources to the storage.
// 	dataSources := []types.DataSource{
// 		types.NewDataSource(
// 			Alice.Address, "NAME1", "DESCRIPTION1", []byte("executable1"),
// 		),
// 		types.NewDataSource(
// 			Bob.Address, "NAME2", "DESCRIPTION2", []byte("executable2"),
// 		),
// 	}
// 	k.SetDataSource(ctx, 1, dataSources[0])
// 	k.SetDataSource(ctx, 2, dataSources[1])
// 	// We should now be able to get all the existing data sources.
// 	require.Equal(t, dataSources, k.GetAllDataSources(ctx))
// }

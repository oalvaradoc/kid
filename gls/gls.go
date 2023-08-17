package gls

import (
	"context"
	"git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/client/mesh"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/model/glsdef"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"git.multiverse.io/eventkit/kit/handler/remote"
	"git.multiverse.io/eventkit/kit/sed/callback"
	"time"
)

const (
	RecordNotFound = "GLSERR0001"
)

type ShardingDataOperator interface {
	/*For `GlsCreate/GlsCreateDxc` start*/
	BindWithTopicInfo(ctx context.Context, topicInfo *TopicInfo, shardingData *ShardingData, opts ...Option) (*SuInfo, *errors.Error)
	BindWithSUType(ctx context.Context, suType string, shardingData *ShardingData, opts ...Option) (*SuInfo, *errors.Error)
	BindListWithTopicInfo(ctx context.Context, topicInfo *TopicInfo, shardingDatas []*ShardingData, opts ...Option) (*SuInfo, *errors.Error)
	BindListWithSUType(ctx context.Context, suType string, shardingDatas []*ShardingData, opts ...Option) (*SuInfo, *errors.Error)
	/*For `GlsCreate/GlsCreateDxc` end*/

	/* For `GlsUpdate/GlsUpdateDxc` start*/
	UnBindWithTopicInfo(ctx context.Context, topicInfo *TopicInfo, shardingData *ShardingData, opts ...Option) (*SuInfo, *errors.Error)
	UnBindWithSUType(ctx context.Context, suType string, shardingData *ShardingData, opts ...Option) (*SuInfo, *errors.Error)
	UnBindListWithTopicInfo(ctx context.Context, topicInfo *TopicInfo, shardingDatas []*ShardingData, opts ...Option) (*SuInfo, *errors.Error)
	UnBindListWithSUType(ctx context.Context, suType string, shardingDatas []*ShardingData, opts ...Option) (*SuInfo, *errors.Error)

	ReBindWithTopicInfo(ctx context.Context, topicInfo *TopicInfo, rebindShardingData *ReBindShardingData, opts ...Option) (*SuInfo, *errors.Error)
	ReBindWithSUType(ctx context.Context, suType string, rebindShardingData *ReBindShardingData, opts ...Option) (*SuInfo, *errors.Error)
	ReBindListWithTopicInfo(ctx context.Context, topicInfo *TopicInfo, rebindShardingDatas []*ReBindShardingData, opts ...Option) (*SuInfo, *errors.Error)
	ReBindListWithSUType(ctx context.Context, suType string, rebindShardingDatas []*ReBindShardingData, opts ...Option) (*SuInfo, *errors.Error)

	AppendIntoBoundRelationWithTopicInfo(ctx context.Context, topicInfo *TopicInfo, appendTo *ShardingData, newShardingData *ShardingData, opts ...Option) (*SuInfo, *errors.Error)
	AppendIntoBoundRelationWithSUType(ctx context.Context, suType string, appendTo *ShardingData, newShardingData *ShardingData, opts ...Option) (*SuInfo, *errors.Error)
	AppendListIntoBoundRelationWithTopicInfo(ctx context.Context, topicInfo *TopicInfo, appendTo *ShardingData, newShardingDatas []*ShardingData, opts ...Option) (*SuInfo, *errors.Error)
	AppendListIntoBoundRelationWithSUType(ctx context.Context, suType string, appendTo *ShardingData, newShardingDatas []*ShardingData, opts ...Option) (*SuInfo, *errors.Error)

	AppendIntoSUWithTopicInfo(ctx context.Context, topicInfo *TopicInfo, appendToSU string, newShardingData *ShardingData, opts ...Option) (*SuInfo, *errors.Error)
	AppendIntoSUWithSUType(ctx context.Context, suType string, appendToSU string, newShardingData *ShardingData, opts ...Option) (*SuInfo, *errors.Error)
	AppendListIntoSUWithTopicInfo(ctx context.Context, topicInfo *TopicInfo, appendToSU string, newShardingDatas []*ShardingData, opts ...Option) (*SuInfo, *errors.Error)
	AppendListIntoSUWithSUType(ctx context.Context, suType string, appendToSU string, newShardingDatas []*ShardingData, opts ...Option) (*SuInfo, *errors.Error)
	/* For `GlsUpdate/GlsUpdateDxc` end*/

	/* For `GlsExist` start */
	IsBoundWithTopicInfo(ctx context.Context, topicInfo *TopicInfo, shardingData *ShardingData, opts ...Option) (bool, *errors.Error)
	IsBoundWithSUType(ctx context.Context, suType string, shardingData *ShardingData, opts ...Option) (bool, *errors.Error)
	/* For `GlsExist` end */

	/* For `GlsRemove/GlsRemoveDxc` start */
	RemoveAllBoundRelation(ctx context.Context, shardingData *ShardingData, opts ...Option) *errors.Error
	/* For `GlsRemove/GlsRemoveDxc` end */

	/* For `GlsSuList` start */
	QuerySuListUsingTopicInfo(ctx context.Context, topicInfo *TopicInfo, opts ...Option) (*SuTypeInfo, *errors.Error)
	QuerySuListUsingSUType(ctx context.Context, suType string, opts ...Option) (*SuTypeInfo, *errors.Error)
	QuerySuListUsingTopicInfoList(ctx context.Context, topicInfoList []TopicInfo, opts ...Option) ([]*SuTypeInfo, *errors.Error)
	QuerySuListUsingSUTypeList(ctx context.Context, suTypeList []string, opts ...Option) ([]*SuTypeInfo, *errors.Error)
	/* For `GlsSuList` end */

	/* For `Lookup` start */
	LookupUsingTopicInfo(ctx context.Context, topicInfo *TopicInfo, shardingData *ShardingData, opts ...Option) (*SuInfo, *errors.Error)
	LookupUsingSUType(ctx context.Context, suType string, shardingData *ShardingData, opts ...Option) (*SuInfo, *errors.Error)
	/* For `Lookup` end */

	/* For `Lookups` start */
	LookupListUsingTopicInfo(ctx context.Context, topicInfo *TopicInfo, shardingDatas []*ShardingData, opts ...Option) ([]*SuInfo, *errors.Error)
	LookupListUsingSUType(ctx context.Context, suType string, shardingDatas []*ShardingData, opts ...Option) ([]*SuInfo, *errors.Error)
	/* For `Lookups` end */
}

type defaultShardingDataOperator struct {
	client       remote.CallInc
	dstGlsSu     string
	glsInterface glsdef.GLSInterface
}

func NewOperator(client remote.CallInc, dstGlsSu string) ShardingDataOperator {
	return &defaultShardingDataOperator{
		client:       client,
		dstGlsSu:     dstGlsSu,
		glsInterface: &callback.GlsOperate{},
	}
}

func NewOptions() Options {
	return Options{
		Timeout:                                 30 * time.Second,
		MaxWaitingTime:                          30 * time.Second,
		RetryWaitingMilliseconds:                100 * time.Millisecond,
		Version:                                 DefaultVersion,
		MaxRetryTimes:                           0,
		DeleteTransactionPropagationInformation: false,
	}
}

// Defines all element status
const (
	GlsUpdateElementInsState    = 1 /*  Update GlsData State: Create, you can ignore this value */
	GlsUpdateElementUpdState    = 2 /*  Update GlsData State: Update, you can ignore this value */
	GlsUpdateElementDelState    = 9 /*  Update GlsData State: Delete, if you want to delete GlsData, you should use this value */
	GlsUpdateElementIgnoreState = 0 /*  Update GlsData State: default 0 is Create */

	DefaultVersion                = "v1"
	BindShardingDataTopicID       = "GlsCreateDxc"
	UpdateShardingDataTopicID     = "GlsUpdateDxc"
	ShardingDataExistCheckTopicID = "GlsExist"
	ShardingDataRemoveTopicID     = "GlsRemoveDxc"
	QuerySUListTopicID            = "GlsSuList"
)

func (d *defaultShardingDataOperator) syncCallAdapter(
	ctx context.Context, options Options, request client.Request, response interface{}, opts ...client.CallOption) (
	client.ResponseMeta, *errors.Error) {
	request.WithOptions(
		mesh.WithTopicTypeBusiness(),
		mesh.WithSU(d.dstGlsSu),
		mesh.WithEventID(options.TargetEventID),
		mesh.WithTimeout(options.Timeout),
		mesh.WithMaxWaitingTime(options.MaxWaitingTime),
		mesh.WithRetryWaitingMilliseconds(options.RetryWaitingMilliseconds),
		mesh.WithVersion(options.Version),
		mesh.WithMaxRetryTimes(options.MaxRetryTimes),
		mesh.WithDeleteTransactionPropagationInformation(options.DeleteTransactionPropagationInformation),
	)
	// sync call
	handlerContext := contexts.BuildHandlerContexts(contexts.ResponseAutoParseKeyMapping(map[string]string{
		"type":            "json",
		"errorCodeKey":    "code",
		"errorMsgKey":     "message",
		"responseDataKey": "response",
	}))
	newCtx := contexts.BuildContextFromParentWithHandlerContexts(ctx, handlerContext)
	responseMeta, err := d.client.SyncCalls(newCtx, request, response)
	if nil != err {
		return nil, errors.Errorf(constant.SystemInternalError, "Failed to sync call, error:%++v", err)
	}

	return responseMeta, nil
}

func (d *defaultShardingDataOperator) BindWithTopicInfo(
	ctx context.Context, topicInfo *TopicInfo, shardingData *ShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = BindShardingDataTopicID
	args := &CreateShardingDataRequest{}
	for _, o := range opts {
		o(&options)
	}

	args.RefragmentShardingDataID = options.RefragmentShardingDataID
	dimension := &DataDimension{
		Topic: GlsTopic{
			TopicType: topicInfo.Type,
			TopicID:   topicInfo.ID,
		},
	}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension

	args.Elements = []GlsData{GlsData{
		GlsType:  shardingData.Type,
		GlsClass: shardingData.Class,
		DataID:   shardingData.ID,
	}}

	request := mesh.NewMeshRequest(args)
	// set response struct
	response := &GlsPrimarySu{}

	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	return &SuInfo{
		Type: response.SuType,
		ID:   response.SuID,
	}, nil
}

func (d *defaultShardingDataOperator) BindWithSUType(
	ctx context.Context, suType string, shardingData *ShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = BindShardingDataTopicID
	args := &CreateShardingDataRequest{}
	for _, o := range opts {
		o(&options)
	}

	args.RefragmentShardingDataID = options.RefragmentShardingDataID
	dimension := &DataDimension{
		SuType: suType,
	}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension

	args.Elements = []GlsData{GlsData{
		GlsType:  shardingData.Type,
		GlsClass: shardingData.Class,
		DataID:   shardingData.ID,
	}}
	request := mesh.NewMeshRequest(args)
	// set response struct
	response := &GlsPrimarySu{}
	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	return &SuInfo{
		Type: response.SuType,
		ID:   response.SuID,
	}, nil
}
func (d *defaultShardingDataOperator) BindListWithTopicInfo(
	ctx context.Context, topicInfo *TopicInfo, shardingDatas []*ShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = BindShardingDataTopicID
	args := &CreateShardingDataRequest{}
	for _, o := range opts {
		o(&options)
	}

	args.RefragmentShardingDataID = options.RefragmentShardingDataID
	dimension := &DataDimension{
		Topic: GlsTopic{
			TopicType: topicInfo.Type,
			TopicID:   topicInfo.ID,
		},
	}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension

	elements := make([]GlsData, 0)
	for _, shardingData := range shardingDatas {
		elements = append(elements, GlsData{
			GlsType:  shardingData.Type,
			GlsClass: shardingData.Class,
			DataID:   shardingData.ID,
		})
	}
	args.Elements = elements
	request := mesh.NewMeshRequest(args)
	// set response struct
	response := &GlsPrimarySu{}
	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	return &SuInfo{
		Type: response.SuType,
		ID:   response.SuID,
	}, nil
}
func (d *defaultShardingDataOperator) BindListWithSUType(
	ctx context.Context, suType string, shardingDatas []*ShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = BindShardingDataTopicID
	args := &CreateShardingDataRequest{}
	for _, o := range opts {
		o(&options)
	}

	args.RefragmentShardingDataID = options.RefragmentShardingDataID
	dimension := &DataDimension{
		SuType: suType,
	}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension

	elements := make([]GlsData, 0)
	for _, shardingData := range shardingDatas {
		elements = append(elements, GlsData{
			GlsType:  shardingData.Type,
			GlsClass: shardingData.Class,
			DataID:   shardingData.ID,
		})
	}
	args.Elements = elements
	request := mesh.NewMeshRequest(args)
	// set response struct
	response := &GlsPrimarySu{}

	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	return &SuInfo{
		Type: response.SuType,
		ID:   response.SuID,
	}, nil
}

func (d *defaultShardingDataOperator) UnBindWithTopicInfo(
	ctx context.Context, topicInfo *TopicInfo, shardingData *ShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = UpdateShardingDataTopicID
	args := &GlsUpdatePrimaryParameter{}
	for _, o := range opts {
		o(&options)
	}

	dimension := &DataDimension{
		Topic: GlsTopic{
			TopicType: topicInfo.Type,
			TopicID:   topicInfo.ID,
		},
	}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension
	args.Primary = GlsData{
		GlsType:  shardingData.Type,
		GlsClass: shardingData.Class,
		DataID:   shardingData.ID,
	}
	args.Elements = []GlsUpdateElement{GlsUpdateElement{
		GlsData: GlsData{
			GlsType:  shardingData.Type,
			GlsClass: shardingData.Class,
			DataID:   shardingData.ID,
		},
		SourceID:    shardingData.ID,
		SourceClass: shardingData.Class,
		State:       GlsUpdateElementDelState,
	}}

	request := mesh.NewMeshRequest(args)
	// set response struct
	response := &GlsPrimarySu{}

	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	return &SuInfo{
		Type: response.SuType,
		ID:   response.SuID,
	}, nil
}
func (d *defaultShardingDataOperator) UnBindWithSUType(
	ctx context.Context, suType string, shardingData *ShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = UpdateShardingDataTopicID
	args := &GlsUpdatePrimaryParameter{}
	for _, o := range opts {
		o(&options)
	}

	dimension := &DataDimension{
		SuType: suType,
	}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension
	args.Primary = GlsData{
		GlsType:  shardingData.Type,
		GlsClass: shardingData.Class,
		DataID:   shardingData.ID,
	}
	args.Elements = []GlsUpdateElement{GlsUpdateElement{
		GlsData: GlsData{
			GlsType:  shardingData.Type,
			GlsClass: shardingData.Class,
			DataID:   shardingData.ID,
		},
		SourceID:    shardingData.ID,
		SourceClass: shardingData.Class,
		State:       GlsUpdateElementDelState,
	}}

	request := mesh.NewMeshRequest(args)
	// set response struct
	response := &GlsPrimarySu{}

	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	return &SuInfo{
		Type: response.SuType,
		ID:   response.SuID,
	}, nil
}
func (d *defaultShardingDataOperator) UnBindListWithTopicInfo(
	ctx context.Context, topicInfo *TopicInfo, shardingDatas []*ShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = UpdateShardingDataTopicID
	args := &GlsUpdatePrimaryParameter{}
	for _, o := range opts {
		o(&options)
	}

	dimension := &DataDimension{
		Topic: GlsTopic{
			TopicType: topicInfo.Type,
			TopicID:   topicInfo.ID,
		},
	}
	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension

	glsUpdateElements := make([]GlsUpdateElement, 0)

	for index, shardingData := range shardingDatas {
		if 0 == index {
			args.Primary = GlsData{
				GlsType:  shardingData.Type,
				GlsClass: shardingData.Class,
				DataID:   shardingData.ID,
			}
		}
		glsUpdateElements = append(glsUpdateElements, GlsUpdateElement{
			GlsData: GlsData{
				GlsType:  shardingData.Type,
				GlsClass: shardingData.Class,
				DataID:   shardingData.ID,
			},
			SourceID:    shardingData.ID,
			SourceClass: shardingData.Class,
			State:       GlsUpdateElementDelState,
		})
	}

	args.Elements = glsUpdateElements

	request := mesh.NewMeshRequest(args)
	// set response struct
	response := &GlsPrimarySu{}

	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	return &SuInfo{
		Type: response.SuType,
		ID:   response.SuID,
	}, nil
}
func (d *defaultShardingDataOperator) UnBindListWithSUType(
	ctx context.Context, suType string, shardingDatas []*ShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = UpdateShardingDataTopicID
	args := &GlsUpdatePrimaryParameter{}
	for _, o := range opts {
		o(&options)
	}

	dimension := &DataDimension{
		SuType: suType,
	}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension

	glsUpdateElements := make([]GlsUpdateElement, 0)

	for index, shardingData := range shardingDatas {
		if 0 == index {
			args.Primary = GlsData{
				GlsType:  shardingData.Type,
				GlsClass: shardingData.Class,
				DataID:   shardingData.ID,
			}
		}
		glsUpdateElements = append(glsUpdateElements, GlsUpdateElement{
			GlsData: GlsData{
				GlsType:  shardingData.Type,
				GlsClass: shardingData.Class,
				DataID:   shardingData.ID,
			},
			SourceID:    shardingData.ID,
			SourceClass: shardingData.Class,
			State:       GlsUpdateElementDelState,
		})
	}

	args.Elements = glsUpdateElements

	request := mesh.NewMeshRequest(args)
	// set response struct
	response := &GlsPrimarySu{}

	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	return &SuInfo{
		Type: response.SuType,
		ID:   response.SuID,
	}, nil
}

func (d *defaultShardingDataOperator) ReBindWithTopicInfo(ctx context.Context, topicInfo *TopicInfo, rebindShardingData *ReBindShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	justRebindShardingDataID := rebindShardingData.Source.Class == rebindShardingData.Target.Class &&
		rebindShardingData.Source.Type == rebindShardingData.Target.Type

	options := NewOptions()
	options.TargetEventID = UpdateShardingDataTopicID
	args := &GlsUpdatePrimaryParameter{}
	for _, o := range opts {
		o(&options)
	}

	dimension := &DataDimension{
		Topic: GlsTopic{
			TopicType: topicInfo.Type,
			TopicID:   topicInfo.ID,
		},
	}
	args.Primary = GlsData{
		GlsType:  rebindShardingData.Source.Type,
		GlsClass: rebindShardingData.Source.Class,
		DataID:   rebindShardingData.Source.ID,
	}
	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension

	if justRebindShardingDataID {
		// if just rebind sharding data id, then just call API one time.
		args.Elements = []GlsUpdateElement{GlsUpdateElement{
			GlsData: GlsData{
				GlsType:  rebindShardingData.Source.Type,
				GlsClass: rebindShardingData.Source.Class,
				DataID:   rebindShardingData.Target.ID,
			},
			SourceID:    rebindShardingData.Source.ID,
			SourceClass: rebindShardingData.Source.Class,
			State:       GlsUpdateElementUpdState,
		}}

		request := mesh.NewMeshRequest(args)
		// set response struct
		response := &GlsPrimarySu{}

		_, err := d.syncCallAdapter(ctx, options, request, response)
		if nil != err {
			return nil, err
		}

		return &SuInfo{
			Type: response.SuType,
			ID:   response.SuID,
		}, nil
	} else {
		// if rebind sharding data new gls type/class then need call `GlsUpdateDxc` and `GlsCreateDxc`
		// delete the source sharding data
		args.Elements = []GlsUpdateElement{GlsUpdateElement{
			GlsData: GlsData{
				GlsType:  rebindShardingData.Source.Type,
				GlsClass: rebindShardingData.Source.Class,
				DataID:   rebindShardingData.Source.ID,
			},
			State: GlsUpdateElementDelState,
		}}

		request := mesh.NewMeshRequest(args)
		// set response struct
		response := &GlsPrimarySu{}

		_, err := d.syncCallAdapter(ctx, options, request, response)
		if nil != err {
			return nil, err
		}

		// create a new sharding data
		options.TargetEventID = BindShardingDataTopicID
		args2 := &CreateShardingDataRequest{}

		args2.RefragmentShardingDataID = options.RefragmentShardingDataID
		dimension2 := &DataDimension{
			Topic: GlsTopic{
				TopicType: topicInfo.Type,
				TopicID:   topicInfo.ID,
			},
		}

		optionalDimension2 := options.OptionalDimension
		if nil != optionalDimension2 {
			dimension2.Tenant = optionalDimension2.Organization
			dimension2.Workspace = optionalDimension2.Workspace
			dimension2.Environment = optionalDimension2.Environment
		}
		args2.Dimension = dimension2

		args2.Elements = []GlsData{GlsData{
			GlsType:  rebindShardingData.Target.Type,
			GlsClass: rebindShardingData.Target.Class,
			DataID:   rebindShardingData.Target.ID,
		}}

		request2 := mesh.NewMeshRequest(args2)
		// set response struct
		response2 := &GlsPrimarySu{}

		_, err = d.syncCallAdapter(ctx, options, request2, response2)
		if nil != err {
			return nil, err
		}

		return &SuInfo{
			Type: response2.SuType,
			ID:   response2.SuID,
		}, nil
	}
}
func (d *defaultShardingDataOperator) ReBindWithSUType(ctx context.Context, suType string, rebindShardingData *ReBindShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	justRebindShardingDataID := rebindShardingData.Source.Class == rebindShardingData.Target.Class &&
		rebindShardingData.Source.Type == rebindShardingData.Target.Type

	options := NewOptions()
	options.TargetEventID = UpdateShardingDataTopicID
	args := &GlsUpdatePrimaryParameter{}
	for _, o := range opts {
		o(&options)
	}

	dimension := &DataDimension{
		SuType: suType,
	}
	args.Primary = GlsData{
		GlsType:  rebindShardingData.Source.Type,
		GlsClass: rebindShardingData.Source.Class,
		DataID:   rebindShardingData.Source.ID,
	}
	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension

	if justRebindShardingDataID {
		// if just rebind sharding data id, then just call API one time.
		args.Elements = []GlsUpdateElement{GlsUpdateElement{
			GlsData: GlsData{
				GlsType:  rebindShardingData.Source.Type,
				GlsClass: rebindShardingData.Source.Class,
				DataID:   rebindShardingData.Target.ID,
			},
			SourceID:    rebindShardingData.Source.ID,
			SourceClass: rebindShardingData.Source.Class,
			State:       GlsUpdateElementUpdState,
		}}

		request := mesh.NewMeshRequest(args)
		// set response struct
		response := &GlsPrimarySu{}

		_, err := d.syncCallAdapter(ctx, options, request, response)
		if nil != err {
			return nil, err
		}

		return &SuInfo{
			Type: response.SuType,
			ID:   response.SuID,
		}, nil
	} else {
		// if rebind sharding data new gls type/class then need call `GlsUpdateDxc` and `GlsCreateDxc`
		// delete the source sharding data
		args.Elements = []GlsUpdateElement{GlsUpdateElement{
			GlsData: GlsData{
				GlsType:  rebindShardingData.Source.Type,
				GlsClass: rebindShardingData.Source.Class,
				DataID:   rebindShardingData.Source.ID,
			},
			State: GlsUpdateElementDelState,
		}}

		request := mesh.NewMeshRequest(args)
		// set response struct
		response := &GlsPrimarySu{}

		_, err := d.syncCallAdapter(ctx, options, request, response)
		if nil != err {
			return nil, err
		}

		// create a new sharding data
		options.TargetEventID = BindShardingDataTopicID
		args2 := &CreateShardingDataRequest{}

		args2.RefragmentShardingDataID = options.RefragmentShardingDataID
		dimension2 := &DataDimension{
			SuType: suType,
		}

		optionalDimension2 := options.OptionalDimension
		if nil != optionalDimension2 {
			dimension2.Tenant = optionalDimension2.Organization
			dimension2.Workspace = optionalDimension2.Workspace
			dimension2.Environment = optionalDimension2.Environment
		}
		args2.Dimension = dimension2

		args2.Elements = []GlsData{GlsData{
			GlsType:  rebindShardingData.Target.Type,
			GlsClass: rebindShardingData.Target.Class,
			DataID:   rebindShardingData.Target.ID,
		}}

		request2 := mesh.NewMeshRequest(args2)
		// set response struct
		response2 := &GlsPrimarySu{}

		_, err = d.syncCallAdapter(ctx, options, request2, response2)
		if nil != err {
			return nil, err
		}

		return &SuInfo{
			Type: response2.SuType,
			ID:   response2.SuID,
		}, nil
	}
}
func (d *defaultShardingDataOperator) ReBindListWithTopicInfo(ctx context.Context, topicInfo *TopicInfo, rebindShardingDatas []*ReBindShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = UpdateShardingDataTopicID
	args := &GlsUpdatePrimaryParameter{}
	for _, o := range opts {
		o(&options)
	}

	dimension := &DataDimension{
		Topic: GlsTopic{
			TopicType: topicInfo.Type,
			TopicID:   topicInfo.ID,
		},
	}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension

	// if rebind sharding data new gls type/class then need call `GlsUpdateDxc` and `GlsCreateDxc`
	// delete the source sharding data
	glsDatas1 := make([]GlsUpdateElement, 0)
	for index, rebindShardingData := range rebindShardingDatas {
		if 0 == index {
			args.Primary = GlsData{
				GlsType:  rebindShardingData.Source.Type,
				GlsClass: rebindShardingData.Source.Class,
				DataID:   rebindShardingData.Source.ID,
			}
		}
		glsDatas1 = append(glsDatas1, GlsUpdateElement{
			GlsData: GlsData{
				GlsType:  rebindShardingData.Source.Type,
				GlsClass: rebindShardingData.Source.Class,
				DataID:   rebindShardingData.Source.ID,
			},
			State: GlsUpdateElementDelState,
		})
	}
	args.Elements = glsDatas1
	request := mesh.NewMeshRequest(args)
	// set response struct
	response := &GlsPrimarySu{}

	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	// create a new sharding data
	options.TargetEventID = BindShardingDataTopicID
	args2 := &CreateShardingDataRequest{}

	args2.RefragmentShardingDataID = options.RefragmentShardingDataID
	dimension2 := &DataDimension{
		Topic: GlsTopic{
			TopicType: topicInfo.Type,
			TopicID:   topicInfo.ID,
		},
	}

	optionalDimension2 := options.OptionalDimension
	if nil != optionalDimension2 {
		dimension2.Tenant = optionalDimension2.Organization
		dimension2.Workspace = optionalDimension2.Workspace
		dimension2.Environment = optionalDimension2.Environment
	}
	args2.Dimension = dimension2

	glsDatas2 := make([]GlsData, 0)
	for _, rebindShardingData := range rebindShardingDatas {
		glsDatas2 = append(glsDatas2, GlsData{
			GlsType:  rebindShardingData.Target.Type,
			GlsClass: rebindShardingData.Target.Class,
			DataID:   rebindShardingData.Target.ID,
		})
	}
	args2.Elements = glsDatas2

	request2 := mesh.NewMeshRequest(args2)
	// set response struct
	response2 := &GlsPrimarySu{}

	_, err = d.syncCallAdapter(ctx, options, request2, response2)
	if nil != err {
		return nil, err
	}

	return &SuInfo{
		Type: response2.SuType,
		ID:   response2.SuID,
	}, nil
}
func (d *defaultShardingDataOperator) ReBindListWithSUType(ctx context.Context, suType string, rebindShardingDatas []*ReBindShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = UpdateShardingDataTopicID
	args := &GlsUpdatePrimaryParameter{}
	for _, o := range opts {
		o(&options)
	}

	dimension := &DataDimension{
		SuType: suType,
	}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension

	// if rebind sharding data new gls type/class then need call `GlsUpdateDxc` and `GlsCreateDxc`
	// delete the source sharding data
	glsDatas1 := make([]GlsUpdateElement, 0)
	for index, rebindShardingData := range rebindShardingDatas {
		if 0 == index {
			args.Primary = GlsData{
				GlsType:  rebindShardingData.Source.Type,
				GlsClass: rebindShardingData.Source.Class,
				DataID:   rebindShardingData.Source.ID,
			}
		}
		glsDatas1 = append(glsDatas1, GlsUpdateElement{
			GlsData: GlsData{
				GlsType:  rebindShardingData.Source.Type,
				GlsClass: rebindShardingData.Source.Class,
				DataID:   rebindShardingData.Source.ID,
			},
			State: GlsUpdateElementDelState,
		})
	}
	args.Elements = glsDatas1
	request := mesh.NewMeshRequest(args)
	// set response struct
	response := &GlsPrimarySu{}

	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	// create a new sharding data
	options.TargetEventID = BindShardingDataTopicID
	args2 := &CreateShardingDataRequest{}

	args2.RefragmentShardingDataID = options.RefragmentShardingDataID
	dimension2 := &DataDimension{
		SuType: suType,
	}

	optionalDimension2 := options.OptionalDimension
	if nil != optionalDimension2 {
		dimension2.Tenant = optionalDimension2.Organization
		dimension2.Workspace = optionalDimension2.Workspace
		dimension2.Environment = optionalDimension2.Environment
	}
	args2.Dimension = dimension2

	glsDatas2 := make([]GlsData, 0)
	for _, rebindShardingData := range rebindShardingDatas {
		glsDatas2 = append(glsDatas2, GlsData{
			GlsType:  rebindShardingData.Target.Type,
			GlsClass: rebindShardingData.Target.Class,
			DataID:   rebindShardingData.Target.ID,
		})
	}
	args2.Elements = glsDatas2

	request2 := mesh.NewMeshRequest(args2)
	// set response struct
	response2 := &GlsPrimarySu{}

	_, err = d.syncCallAdapter(ctx, options, request2, response2)
	if nil != err {
		return nil, err
	}

	return &SuInfo{
		Type: response2.SuType,
		ID:   response2.SuID,
	}, nil
}

func (d *defaultShardingDataOperator) AppendIntoBoundRelationWithTopicInfo(ctx context.Context, topicInfo *TopicInfo, appendTo *ShardingData, newShardingData *ShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = UpdateShardingDataTopicID
	args := &GlsUpdatePrimaryParameter{}
	for _, o := range opts {
		o(&options)
	}

	dimension := &DataDimension{
		Topic: GlsTopic{
			TopicType: topicInfo.Type,
			TopicID:   topicInfo.ID,
		},
	}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension
	args.RefragmentShardingDataID = options.RefragmentShardingDataID
	args.Primary = GlsData{
		GlsType:  appendTo.Type,
		GlsClass: appendTo.Class,
		DataID:   appendTo.ID,
	}
	args.Elements = []GlsUpdateElement{GlsUpdateElement{
		GlsData: GlsData{
			GlsType:  newShardingData.Type,
			GlsClass: newShardingData.Class,
			DataID:   newShardingData.ID,
		},
		State: GlsUpdateElementIgnoreState,
	}}

	request := mesh.NewMeshRequest(args)
	// set response struct
	response := &GlsPrimarySu{}

	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	return &SuInfo{
		Type: response.SuType,
		ID:   response.SuID,
	}, nil
}
func (d *defaultShardingDataOperator) AppendIntoBoundRelationWithSUType(ctx context.Context, suType string, appendTo *ShardingData, newShardingData *ShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = UpdateShardingDataTopicID
	args := &GlsUpdatePrimaryParameter{}
	for _, o := range opts {
		o(&options)
	}

	dimension := &DataDimension{
		SuType: suType,
	}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension
	args.RefragmentShardingDataID = options.RefragmentShardingDataID
	args.Primary = GlsData{
		GlsType:  appendTo.Type,
		GlsClass: appendTo.Class,
		DataID:   appendTo.ID,
	}
	args.Elements = []GlsUpdateElement{GlsUpdateElement{
		GlsData: GlsData{
			GlsType:  newShardingData.Type,
			GlsClass: newShardingData.Class,
			DataID:   newShardingData.ID,
		},
		State: GlsUpdateElementIgnoreState,
	}}

	request := mesh.NewMeshRequest(args)
	// set response struct
	response := &GlsPrimarySu{}

	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	return &SuInfo{
		Type: response.SuType,
		ID:   response.SuID,
	}, nil
}
func (d *defaultShardingDataOperator) AppendListIntoBoundRelationWithTopicInfo(ctx context.Context, topicInfo *TopicInfo, appendTo *ShardingData, newShardingDatas []*ShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = UpdateShardingDataTopicID
	args := &GlsUpdatePrimaryParameter{}
	for _, o := range opts {
		o(&options)
	}

	dimension := &DataDimension{
		Topic: GlsTopic{
			TopicType: topicInfo.Type,
			TopicID:   topicInfo.ID,
		},
	}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension
	args.RefragmentShardingDataID = options.RefragmentShardingDataID
	args.Primary = GlsData{
		GlsType:  appendTo.Type,
		GlsClass: appendTo.Class,
		DataID:   appendTo.ID,
	}
	glsUpdateElements := make([]GlsUpdateElement, 0)
	for _, newShardingData := range newShardingDatas {
		glsUpdateElements = append(glsUpdateElements, GlsUpdateElement{
			GlsData: GlsData{
				GlsType:  newShardingData.Type,
				GlsClass: newShardingData.Class,
				DataID:   newShardingData.ID,
			},
			State: GlsUpdateElementIgnoreState,
		})
	}
	args.Elements = glsUpdateElements

	request := mesh.NewMeshRequest(args)
	// set response struct
	response := &GlsPrimarySu{}

	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	return &SuInfo{
		Type: response.SuType,
		ID:   response.SuID,
	}, nil
}
func (d *defaultShardingDataOperator) AppendListIntoBoundRelationWithSUType(ctx context.Context, suType string, appendTo *ShardingData, newShardingDatas []*ShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = UpdateShardingDataTopicID
	args := &GlsUpdatePrimaryParameter{}
	for _, o := range opts {
		o(&options)
	}

	dimension := &DataDimension{
		SuType: suType,
	}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension
	args.RefragmentShardingDataID = options.RefragmentShardingDataID
	args.Primary = GlsData{
		GlsType:  appendTo.Type,
		GlsClass: appendTo.Class,
		DataID:   appendTo.ID,
	}
	glsUpdateElements := make([]GlsUpdateElement, 0)
	for _, newShardingData := range newShardingDatas {
		glsUpdateElements = append(glsUpdateElements, GlsUpdateElement{
			GlsData: GlsData{
				GlsType:  newShardingData.Type,
				GlsClass: newShardingData.Class,
				DataID:   newShardingData.ID,
			},
			State: GlsUpdateElementIgnoreState,
		})
	}
	args.Elements = glsUpdateElements

	request := mesh.NewMeshRequest(args)
	// set response struct
	response := &GlsPrimarySu{}

	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	return &SuInfo{
		Type: response.SuType,
		ID:   response.SuID,
	}, nil
}

func (d *defaultShardingDataOperator) AppendIntoSUWithTopicInfo(ctx context.Context, topicInfo *TopicInfo, appendToSU string, newShardingData *ShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = UpdateShardingDataTopicID
	args := &GlsUpdatePrimaryParameter{}
	for _, o := range opts {
		o(&options)
	}

	dimension := &DataDimension{
		Topic: GlsTopic{
			TopicType: topicInfo.Type,
			TopicID:   topicInfo.ID,
		},
	}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension
	args.RefragmentShardingDataID = options.RefragmentShardingDataID
	args.AppendToSU = appendToSU
	args.Elements = []GlsUpdateElement{GlsUpdateElement{
		GlsData: GlsData{
			GlsType:  newShardingData.Type,
			GlsClass: newShardingData.Class,
			DataID:   newShardingData.ID,
		},
		State: GlsUpdateElementIgnoreState,
	}}

	request := mesh.NewMeshRequest(args)
	// set response struct
	response := &GlsPrimarySu{}

	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	return &SuInfo{
		Type: response.SuType,
		ID:   response.SuID,
	}, nil
}
func (d *defaultShardingDataOperator) AppendIntoSUWithSUType(ctx context.Context, suType string, appendToSU string, newShardingData *ShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = UpdateShardingDataTopicID
	args := &GlsUpdatePrimaryParameter{}
	for _, o := range opts {
		o(&options)
	}

	dimension := &DataDimension{
		SuType: suType,
	}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension
	args.RefragmentShardingDataID = options.RefragmentShardingDataID
	args.AppendToSU = appendToSU
	args.Elements = []GlsUpdateElement{GlsUpdateElement{
		GlsData: GlsData{
			GlsType:  newShardingData.Type,
			GlsClass: newShardingData.Class,
			DataID:   newShardingData.ID,
		},
		State: GlsUpdateElementIgnoreState,
	}}

	request := mesh.NewMeshRequest(args)
	// set response struct
	response := &GlsPrimarySu{}

	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	return &SuInfo{
		Type: response.SuType,
		ID:   response.SuID,
	}, nil
}
func (d *defaultShardingDataOperator) AppendListIntoSUWithTopicInfo(ctx context.Context, topicInfo *TopicInfo, appendToSU string, newShardingDatas []*ShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = UpdateShardingDataTopicID
	args := &GlsUpdatePrimaryParameter{}
	for _, o := range opts {
		o(&options)
	}

	dimension := &DataDimension{
		Topic: GlsTopic{
			TopicType: topicInfo.Type,
			TopicID:   topicInfo.ID,
		},
	}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension
	args.RefragmentShardingDataID = options.RefragmentShardingDataID
	args.AppendToSU = appendToSU
	glsUpdateElements := make([]GlsUpdateElement, 0)
	for _, newShardingData := range newShardingDatas {
		glsUpdateElements = append(glsUpdateElements, GlsUpdateElement{
			GlsData: GlsData{
				GlsType:  newShardingData.Type,
				GlsClass: newShardingData.Class,
				DataID:   newShardingData.ID,
			},
			State: GlsUpdateElementIgnoreState,
		})
	}
	args.Elements = glsUpdateElements

	request := mesh.NewMeshRequest(args)
	// set response struct
	response := &GlsPrimarySu{}

	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	return &SuInfo{
		Type: response.SuType,
		ID:   response.SuID,
	}, nil
}
func (d *defaultShardingDataOperator) AppendListIntoSUWithSUType(ctx context.Context, suType string, appendToSU string, newShardingDatas []*ShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = UpdateShardingDataTopicID
	args := &GlsUpdatePrimaryParameter{}
	for _, o := range opts {
		o(&options)
	}

	dimension := &DataDimension{
		SuType: suType,
	}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension
	args.RefragmentShardingDataID = options.RefragmentShardingDataID
	args.AppendToSU = appendToSU
	glsUpdateElements := make([]GlsUpdateElement, 0)
	for _, newShardingData := range newShardingDatas {
		glsUpdateElements = append(glsUpdateElements, GlsUpdateElement{
			GlsData: GlsData{
				GlsType:  newShardingData.Type,
				GlsClass: newShardingData.Class,
				DataID:   newShardingData.ID,
			},
			State: GlsUpdateElementIgnoreState,
		})
	}
	args.Elements = glsUpdateElements

	request := mesh.NewMeshRequest(args)
	// set response struct
	response := &GlsPrimarySu{}

	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	return &SuInfo{
		Type: response.SuType,
		ID:   response.SuID,
	}, nil
}

/* For `GlsUpdate/GlsUpdateDxc` end*/

/* For `GlsExist` start */
func (d *defaultShardingDataOperator) IsBoundWithTopicInfo(ctx context.Context, topicInfo *TopicInfo, shardingData *ShardingData, opts ...Option) (bool, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = ShardingDataExistCheckTopicID
	args := &ShardingDataExistCheckRequest{}
	for _, o := range opts {
		o(&options)
	}

	dimension := &DataDimension{
		Topic: GlsTopic{
			TopicType: topicInfo.Type,
			TopicID:   topicInfo.ID,
		},
	}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension

	args.Element = GlsData{
		GlsType:  shardingData.Type,
		GlsClass: shardingData.Class,
		DataID:   shardingData.ID,
	}

	request := mesh.NewMeshRequest(args)
	// set response struct
	var response bool

	_, err := d.syncCallAdapter(ctx, options, request, &response)
	if nil != err {
		return false, err
	}

	return response, nil
}
func (d *defaultShardingDataOperator) IsBoundWithSUType(ctx context.Context, suType string, shardingData *ShardingData, opts ...Option) (bool, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = ShardingDataExistCheckTopicID
	args := &ShardingDataExistCheckRequest{}
	for _, o := range opts {
		o(&options)
	}

	dimension := &DataDimension{
		SuType: suType,
	}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		dimension.Tenant = optionalDimension.Organization
		dimension.Workspace = optionalDimension.Workspace
		dimension.Environment = optionalDimension.Environment
	}
	args.Dimension = dimension

	args.Element = GlsData{
		GlsType:  shardingData.Type,
		GlsClass: shardingData.Class,
		DataID:   shardingData.ID,
	}

	request := mesh.NewMeshRequest(args)
	// set response struct
	var response bool

	_, err := d.syncCallAdapter(ctx, options, request, &response)
	if nil != err {
		return false, err
	}

	return response, nil
}

/* For `GlsExist` end */

/* For `GlsRemove/GlsRemoveDxc` start */
func (d *defaultShardingDataOperator) RemoveAllBoundRelation(ctx context.Context, shardingData *ShardingData, opts ...Option) *errors.Error {
	options := NewOptions()
	options.TargetEventID = ShardingDataRemoveTopicID
	args := &GlsData{}
	for _, o := range opts {
		o(&options)
	}
	args.GlsType = shardingData.Type
	args.GlsClass = shardingData.Class
	args.DataID = shardingData.ID

	request := mesh.NewMeshRequest(args)
	_, err := d.syncCallAdapter(ctx, options, request, nil)
	if nil != err {
		return err
	}

	return nil
}

/* For `GlsRemove/GlsRemoveDxc` end */

func (d *defaultShardingDataOperator) QuerySuListUsingTopicInfo(ctx context.Context, topicInfo *TopicInfo, opts ...Option) (*SuTypeInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = QuerySUListTopicID
	for _, o := range opts {
		o(&options)
	}
	args := &QuerySuListRequest{}
	glsDimension := &GlsDimension{}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		glsDimension.Tenant = optionalDimension.Organization
		glsDimension.Workspace = optionalDimension.Workspace
		glsDimension.Environment = optionalDimension.Environment
	}
	args.Dimension = glsDimension
	args.SuTypes = []SuTypeWrapper{
		SuTypeWrapper{
			Topic: GlsTopic{
				TopicType: topicInfo.Type,
				TopicID:   topicInfo.ID,
			},
		},
	}
	request := mesh.NewMeshRequest(args)
	// set response struct

	response := &QuerySuListResponse{}
	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	if len(response.SuTypes) == 0 {
		return nil, errors.Errorf(RecordNotFound, "Cannot found su list using topic info:%++v", topicInfo)
	}

	return &SuTypeInfo{
		TopicInfo: TopicInfo{
			Type: topicInfo.Type,
			ID:   topicInfo.ID,
		},
		SuList: response.SuTypes[0].SuList,
	}, nil
}
func (d *defaultShardingDataOperator) QuerySuListUsingSUType(ctx context.Context, suType string, opts ...Option) (*SuTypeInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = QuerySUListTopicID
	for _, o := range opts {
		o(&options)
	}
	args := &QuerySuListRequest{}
	glsDimension := &GlsDimension{}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		glsDimension.Tenant = optionalDimension.Organization
		glsDimension.Workspace = optionalDimension.Workspace
		glsDimension.Environment = optionalDimension.Environment
	}
	args.Dimension = glsDimension
	args.SuTypes = []SuTypeWrapper{
		SuTypeWrapper{
			SuType: suType,
		},
	}
	request := mesh.NewMeshRequest(args)
	// set response struct
	response := &QuerySuListResponse{}

	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	if len(response.SuTypes) == 0 {
		return nil, errors.Errorf(RecordNotFound, "Cannot found su list using su type:%s", suType)
	}

	return &SuTypeInfo{
		SuType: suType,
		SuList: response.SuTypes[0].SuList,
	}, nil
}

func (d *defaultShardingDataOperator) QuerySuListUsingTopicInfoList(ctx context.Context, topicInfoList []TopicInfo, opts ...Option) ([]*SuTypeInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = QuerySUListTopicID
	for _, o := range opts {
		o(&options)
	}
	args := &QuerySuListRequest{}
	glsDimension := &GlsDimension{}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		glsDimension.Tenant = optionalDimension.Organization
		glsDimension.Workspace = optionalDimension.Workspace
		glsDimension.Environment = optionalDimension.Environment
	}
	args.Dimension = glsDimension

	SuTypeWrapperList := make([]SuTypeWrapper, 0)
	for _, topicInfo := range topicInfoList {
		SuTypeWrapperList = append(SuTypeWrapperList, SuTypeWrapper{
			Topic: GlsTopic{
				TopicType: topicInfo.Type,
				TopicID:   topicInfo.ID,
			},
		})
	}
	args.SuTypes = SuTypeWrapperList

	request := mesh.NewMeshRequest(args)
	// set response struct

	response := &QuerySuListResponse{}
	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	if len(response.SuTypes) == 0 {
		return nil, errors.Errorf(RecordNotFound, "Cannot found su list using topic info list:%++v", topicInfoList)
	}

	suTypeInfoList := make([]*SuTypeInfo, 0)
	for _, suTypeWrapper := range response.SuTypes {
		suTypeInfoList = append(suTypeInfoList, &SuTypeInfo{
			SuType: suTypeWrapper.SuType,
			TopicInfo: TopicInfo{
				Type: suTypeWrapper.Topic.TopicType,
				ID:   suTypeWrapper.Topic.TopicID,
			},
			SuList: suTypeWrapper.SuList,
		})
	}

	return suTypeInfoList, nil
}

func (d *defaultShardingDataOperator) QuerySuListUsingSUTypeList(ctx context.Context, suTypeList []string, opts ...Option) ([]*SuTypeInfo, *errors.Error) {
	options := NewOptions()
	options.TargetEventID = QuerySUListTopicID
	for _, o := range opts {
		o(&options)
	}
	args := &QuerySuListRequest{}
	glsDimension := &GlsDimension{}

	optionalDimension := options.OptionalDimension
	if nil != optionalDimension {
		glsDimension.Tenant = optionalDimension.Organization
		glsDimension.Workspace = optionalDimension.Workspace
		glsDimension.Environment = optionalDimension.Environment
	}
	args.Dimension = glsDimension

	SuTypeWrapperList := make([]SuTypeWrapper, 0)
	for _, suType := range suTypeList {
		SuTypeWrapperList = append(SuTypeWrapperList, SuTypeWrapper{
			SuType: suType,
		})
	}
	args.SuTypes = SuTypeWrapperList
	request := mesh.NewMeshRequest(args)
	// set response struct
	response := &QuerySuListResponse{}

	_, err := d.syncCallAdapter(ctx, options, request, response)
	if nil != err {
		return nil, err
	}

	if len(response.SuTypes) == 0 {
		return nil, errors.Errorf(RecordNotFound, "Cannot found su list using su type list:%s", suTypeList)
	}

	suTypeInfoList := make([]*SuTypeInfo, 0)
	for _, suTypeWrapper := range response.SuTypes {
		suTypeInfoList = append(suTypeInfoList, &SuTypeInfo{
			SuType: suTypeWrapper.SuType,
			TopicInfo: TopicInfo{
				Type: suTypeWrapper.Topic.TopicType,
				ID:   suTypeWrapper.Topic.TopicID,
			},
			SuList: suTypeWrapper.SuList,
		})
	}

	return suTypeInfoList, nil
}

/* For `Lookup` start */
func (d *defaultShardingDataOperator) LookupUsingTopicInfo(ctx context.Context, topicInfo *TopicInfo, shardingData *ShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	dimension := glsdef.Dimension{
		Topic: glsdef.Topic{
			TopicType: topicInfo.Type,
			TopicID:   topicInfo.ID,
		},
	}
	element := &glsdef.Element{
		ElementType:  shardingData.Type,
		ElementClass: shardingData.Class,
		ElementID:    shardingData.ID,
	}
	if element.ElementClass == "" {
		element.ElementClass = element.ElementType
	}

	pd, code, err := d.glsInterface.Lookup(dimension, element)
	if nil != err {
		return nil, errors.Errorf(constant.SystemInternalError, "GLS lookup failed, code[%v] error:[%s]", code, errors.ErrorToString(err))
	}

	return &SuInfo{
		Type: pd.SuType,
		ID:   pd.SuID,
	}, nil
}

func (d *defaultShardingDataOperator) LookupUsingSUType(ctx context.Context, suType string, shardingData *ShardingData, opts ...Option) (*SuInfo, *errors.Error) {
	dimension := glsdef.Dimension{
		SuType: suType,
	}
	element := &glsdef.Element{
		ElementType:  shardingData.Type,
		ElementClass: shardingData.Class,
		ElementID:    shardingData.ID,
	}
	if element.ElementClass == "" {
		element.ElementClass = element.ElementType
	}

	pd, code, err := d.glsInterface.Lookup(dimension, element)
	if nil != err {
		return nil, errors.Errorf(constant.SystemInternalError, "GLS lookup failed, code[%v] error:[%s]", code, errors.ErrorToString(err))
	}

	return &SuInfo{
		Type: pd.SuType,
		ID:   pd.SuID,
	}, nil
}

/* For `Lookup` end */

/* For `Lookups` start */
func (d *defaultShardingDataOperator) LookupListUsingTopicInfo(ctx context.Context, topicInfo *TopicInfo, shardingDatas []*ShardingData, opts ...Option) ([]*SuInfo, *errors.Error) {
	dimension := glsdef.Dimension{
		Topic: glsdef.Topic{
			TopicType: topicInfo.Type,
			TopicID:   topicInfo.ID,
		},
	}
	elements := make([]glsdef.Element, 0)

	for _, shardingData := range shardingDatas {
		element := glsdef.Element{
			ElementType:  shardingData.Type,
			ElementClass: shardingData.Class,
			ElementID:    shardingData.ID,
		}
		if element.ElementClass == "" {
			element.ElementClass = element.ElementType
		}
		elements = append(elements, element)
	}

	pd, code, err := d.glsInterface.Lookups(dimension, elements)
	if nil != err {
		return nil, errors.Errorf(constant.SystemInternalError, "GLS lookup failed, code[%v] error:[%s]", code, errors.ErrorToString(err))
	}
	result := make([]*SuInfo, 0)

	for _, pu := range pd.PrmSus {
		result = append(result, &SuInfo{
			Type: pu.SuType,
			ID:   pu.SuID,
		})
	}
	return result, nil
}

func (d *defaultShardingDataOperator) LookupListUsingSUType(ctx context.Context, suType string, shardingDatas []*ShardingData, opts ...Option) ([]*SuInfo, *errors.Error) {
	dimension := glsdef.Dimension{
		SuType: suType,
	}
	elements := make([]glsdef.Element, 0)

	for _, shardingData := range shardingDatas {
		element := glsdef.Element{
			ElementType:  shardingData.Type,
			ElementClass: shardingData.Class,
			ElementID:    shardingData.ID,
		}
		if element.ElementClass == "" {
			element.ElementClass = element.ElementType
		}
		elements = append(elements, element)
	}

	pd, code, err := d.glsInterface.Lookups(dimension, elements)
	if nil != err {
		return nil, errors.Errorf(constant.SystemInternalError, "GLS lookup failed, code[%v] error:[%s]", code, errors.ErrorToString(err))
	}
	result := make([]*SuInfo, 0)

	for _, pu := range pd.PrmSus {
		result = append(result, &SuInfo{
			Type: pu.SuType,
			ID:   pu.SuID,
		})
	}
	return result, nil
}

/* For `Lookups` end */

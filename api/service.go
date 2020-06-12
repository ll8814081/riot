package api

import (
	"context"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/raft"
	"github.com/laohanlinux/riot/cluster"
)

type NotArg struct{}
type NotReply struct{}

/////////////////////////
type SetKVArg struct {
	BucketName string
	Key        string
	Value      []byte
}

type SetKVReply struct {
	HasBucket bool
}

type GetKVArg struct {
	BucketName string
	Key        string
}

type GetPrefixKVArg struct {
	BucketName string
	KeyPrefix  string
}

type GetKVReply struct {
	Has   bool
	Value []byte
}

type GetPrefixKVReply struct {
	Has bool
	Kv  map[string][]byte
}

type DelKVArg struct {
	BucketName string
	Key        string
}

//////////////

type DelBucketArg struct {
	BucketName string
}

type BucketInfoArg struct {
	BucketName string
}

type CreateBucketArg struct {
	BucketName string
}

type BucketInfoReply struct {
	Has  bool
	Info bolt.BucketStats
}

////////////
type NodeStateReply struct {
	State string
}

type NodeString struct {
	NodeInfo string
}

type PeersReply struct {
	Peers []string
}

type LeaderReply struct {
	Leader string
}

type SnapshotReply struct {
	Len int
}

type RemovePeerArg struct {
	Peer string
}

type RemovePeerReply struct {
	Has bool
}

//////////

type APIService struct {
	api API
	adm AdmAPI
}

func NewAPIService(api API, adm AdmAPI) *APIService {
	return &APIService{api: api, adm: adm}
}

func (s *APIService) KV(_ context.Context, arg *GetKVArg, reply *GetKVReply) (err error) {
	reply.Value, err = s.api.GetValue(arg.BucketName, arg.Key)
	if err == cluster.ErrNotFound {
		err = nil
	} else {
		reply.Has = true
	}
	return
}
func (s *APIService) SetKV(_ context.Context, arg *SetKVArg, reply *SetKVReply) (err error) {
	err = s.api.SetKV(arg.BucketName, arg.Key, arg.Value)
	if err == cluster.ErrNotFound {
		err = nil
		reply.HasBucket = false
	}
	return
}

func (s *APIService) GetPrefixKV(_ context.Context, arg *GetPrefixKVArg, reply *GetPrefixKVReply) (err error) {
	kv, err := s.api.GetPrefixKV(arg.BucketName, arg.KeyPrefix)
	if len(kv) == 0 {
		reply.Has = false
	} else {
		reply.Has = true
	}
	reply.Kv = kv
	return
}

func (s *APIService) BucketInfo(_ context.Context, arg *BucketInfoArg, reply *BucketInfoReply) (err error) {
	var info interface{}
	if info, err = s.api.GetBucket(arg.BucketName); err == cluster.ErrNotFound {
		err = nil
		return
	}
	reply.Info = info.(bolt.BucketStats)
	reply.Has = true
	return
}

func (s *APIService) DelKey(_ context.Context, arg *DelKVArg, _ *NotReply) (err error) {
	if err = s.api.DelKey(arg.BucketName, arg.Key); err == cluster.ErrNotFound {
		err = nil
	}
	return
}

func (s *APIService) DelBucket(_ context.Context, arg *DelBucketArg, _ *NotReply) (err error) {
	if err = s.api.DelBucket(arg.BucketName); err == cluster.ErrNotFound {
		err = nil
	}
	return
}

func (s *APIService) CreateBucket(_ context.Context, arg *CreateBucketArg, _ *NotReply) (err error) {
	err = s.api.CreateBucket(arg.BucketName)
	return
}

func (s *APIService) NodeState(_ context.Context, _ *NotArg, reply *NodeStateReply) (err error) {
	reply.State = s.adm.State()
	return
}

func (s *APIService) NodeString(_ context.Context, _ *NotArg, reply *NodeString) (err error) {
	reply.NodeInfo = s.adm.NodeString()
	return
}

func (s *APIService) Peers(_ context.Context, _ *NotArg, reply *PeersReply) (err error) {
	reply.Peers, err = s.adm.Peers()
	return
}

func (s *APIService) Leader(_ context.Context, _ *NotArg, reply *LeaderReply) (err error) {
	reply.Leader, err = s.adm.Leader()
	return
}

func (s *APIService) Snapshot(_ context.Context, _ *NotArg, reply *SnapshotReply) (err error) {
	reply.Len, err = s.adm.Snapshot()
	return
}

func (s *APIService) RemovePeer(_ context.Context, arg *RemovePeerArg, reply *RemovePeerReply) (err error) {
	if err = s.adm.RemovePeer(arg.Peer); err == raft.ErrUnknownPeer {
		err = nil
	} else {
		reply.Has = true
	}
	return
}

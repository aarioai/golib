package weixintransfer

import (
	"context"
)

// Transfer 发起商家转账
// https://pay.weixin.qq.com/doc/v3/merchant/4012458841
// 注意受理成功将返回批次单号，此时并不代表转账成功，请通过查单接口查询单据的付款状态
func (s *Service) Transfer(ctx context.Context, req TransferRequest) (TransferResult, error) {
	client, err := s.transferClient(ctx)
	if err != nil {
		return TransferResult{}, err
	}

	resp, _, err := client.InitiateBatchTransfer(ctx, req.adapter(s.Appid))
	if err != nil {
		return TransferResult{}, err
	}
	return toTransferResponse(*resp, s.loc), nil
}

// Query 通过商家批次单号查询批次单
// https://pay.weixin.qq.com/doc/v3/merchant/4012458868
func (s *Service) Query(ctx context.Context, req QueryRequest) (QueryResult, error) {
	client, err := s.transferClient(ctx)
	if err != nil {
		return QueryResult{}, err
	}

	resp, _, err := client.GetTransferBatchByOutNo(ctx, req.adapter())
	if err != nil {
		return QueryResult{}, err
	}
	return toQueryResult(*resp, s.loc), nil
}

// QueryByWeixinBatch 通过微信批次单号查询批次单
func (s *Service) QueryByWeixinBatch(ctx context.Context, req QueryRequestByWeixinBatch) (QueryResult, error) {
	client, err := s.transferClient(ctx)
	if err != nil {
		return QueryResult{}, err
	}

	resp, _, err := client.GetTransferBatchByNo(ctx, req.adapter())
	if err != nil {
		return QueryResult{}, err
	}
	return toQueryResult(*resp, s.loc), nil
}

// QueryDetail 通过商家明细单号查询明细单
func (s *Service) QueryDetail(ctx context.Context, req DetailRequest) (DetailResult, error) {
	client, err := s.detailClient(ctx)
	if err != nil {
		return DetailResult{}, err
	}
	resp, _, err := client.GetTransferDetailByOutNo(ctx, req.adapter())
	if err != nil {
		return DetailResult{}, err
	}
	return toDetailResult(*resp, s.loc), nil
}

// QueryDetailByWeixinBatch 通过微信明细单号查询明细单
func (s *Service) QueryDetailByWeixinBatch(ctx context.Context, req DetailRequestByWeixinBatch) (DetailResult, error) {
	client, err := s.detailClient(ctx)
	if err != nil {
		return DetailResult{}, err
	}
	resp, _, err := client.GetTransferDetailByNo(ctx, req.adapter())
	if err != nil {
		return DetailResult{}, err
	}
	return toDetailResult(*resp, s.loc), nil
}

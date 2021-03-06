/*
 * Copyright (c) 2018. Abstrium SAS <team (at) pydio.com>
 * This file is part of Pydio Cells.
 *
 * Pydio Cells is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Pydio Cells is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Pydio Cells.  If not, see <http://www.gnu.org/licenses/>.
 *
 * The latest code can be found at <https://pydio.com>.
 */

package grpc

import (
	"context"

	"github.com/micro/go-micro/client"
	"go.uber.org/zap"

	"github.com/pydio/cells/common"
	"github.com/pydio/cells/common/auth"
	"github.com/pydio/cells/common/log"
	"github.com/pydio/cells/common/proto/idm"
	"github.com/pydio/cells/common/proto/tree"
	"github.com/pydio/cells/common/service/context"
	"github.com/pydio/cells/common/service/proto"
	"github.com/pydio/cells/idm/meta"
)

// Handler definition.
type Handler struct{}

// UpdateUserMeta adds, updates or deletes user meta.
func (h *Handler) UpdateUserMeta(ctx context.Context, request *idm.UpdateUserMetaRequest, response *idm.UpdateUserMetaResponse) error {

	dao := servicecontext.GetDAO(ctx).(meta.DAO)

	namespaces, _ := dao.GetNamespaceDao().List()
	indexableMetas := make(map[string]map[string]string)
	for _, metadata := range request.MetaDatas {
		if request.Operation == idm.UpdateUserMetaRequest_PUT {
			// ADD / UPDATE
			if newMeta, _, err := dao.Set(metadata); err == nil {
				if ns, ok := namespaces[metadata.Namespace]; ok && ns.Indexable {
					var nodesIndexes map[string]string
					var has bool
					if nodesIndexes, has = indexableMetas[metadata.NodeUuid]; !has {
						nodesIndexes = make(map[string]string)
						indexableMetas[metadata.NodeUuid] = nodesIndexes
					}
					nodesIndexes[metadata.Namespace] = metadata.JsonValue
				}
				response.MetaDatas = append(response.MetaDatas, newMeta)
			} else {
				return err
			}
		} else {
			// DELETE
			if err := dao.Del(metadata); err != nil {
				return err
			} else {
				if ns, ok := namespaces[metadata.Namespace]; ok && ns.Indexable {
					var nodesIndexes map[string]string
					var has bool
					if nodesIndexes, has = indexableMetas[metadata.NodeUuid]; !has {
						nodesIndexes = make(map[string]string)
						indexableMetas[metadata.NodeUuid] = nodesIndexes
					}
					nodesIndexes[metadata.Namespace] = ""
				}
			}
		}
	}

	for nodeId, toIndex := range indexableMetas {
		node := &tree.Node{Uuid: nodeId}
		node.MetaStore = toIndex
		log.Logger(ctx).Info("Publishing UPDATE META for node, shall we update the node, or switch to UPDATE_META_DELTA?", node.Zap())
		client.Publish(ctx, client.NewPublication(common.TOPIC_META_CHANGES, &tree.NodeChangeEvent{
			Type:   tree.NodeChangeEvent_UPDATE_META,
			Target: node,
		}))
	}

	return nil

}

// SearchUserMeta retrieves meta based on various criteria.
func (h *Handler) SearchUserMeta(ctx context.Context, request *idm.SearchUserMetaRequest, stream idm.UserMetaService_SearchUserMetaStream) error {

	defer stream.Close()
	dao := servicecontext.GetDAO(ctx).(meta.DAO)
	results, err := dao.Search(request.MetaUuids, request.NodeUuids, request.Namespace, request.ResourceSubjectOwner, request.ResourceQuery)
	if err != nil {
		return err
	}
	for _, result := range results {
		stream.Send(&idm.SearchUserMetaResponse{UserMeta: result})
	}
	return nil

}

// Implements ReadNodeStream to be a meta provider.
func (h *Handler) ReadNodeStream(ctx context.Context, stream tree.NodeProviderStreamer_ReadNodeStreamStream) error {

	defer stream.Close()
	dao := servicecontext.GetDAO(ctx).(meta.DAO)
	subjects, e := auth.SubjectsForResourcePolicyQuery(ctx, nil)
	if e != nil {
		return e
	}

	for {
		req, er := stream.Recv()
		if req == nil {
			break
		}
		if er != nil {
			return er
		}
		node := req.Node

		results, err := dao.Search([]string{}, []string{node.Uuid}, "", "", &service.ResourcePolicyQuery{
			Subjects: subjects,
		})
		log.Logger(ctx).Debug("Got Results For Node", node.ZapUuid(), zap.Any("results", results))
		if err == nil && len(results) > 0 {
			for _, result := range results {
				node.MetaStore[result.Namespace] = result.JsonValue
			}
		}
		stream.Send(&tree.ReadNodeResponse{Node: node})
	}

	return nil
}

// Update/Delete a namespace.
func (h *Handler) UpdateUserMetaNamespace(ctx context.Context, request *idm.UpdateUserMetaNamespaceRequest, response *idm.UpdateUserMetaNamespaceResponse) error {

	dao := servicecontext.GetDAO(ctx).(meta.DAO).GetNamespaceDao()
	for _, metaNameSpace := range request.Namespaces {
		if err := dao.Del(metaNameSpace); err != nil {
			return err
		}
	}
	if request.Operation == idm.UpdateUserMetaNamespaceRequest_PUT {
		for _, metaNameSpace := range request.Namespaces {
			if err := dao.Add(metaNameSpace); err != nil {
				return err
			}
			response.Namespaces = append(response.Namespaces, metaNameSpace)
		}
	}
	return nil

}

// List all namespaces from underlying DAO.
func (h *Handler) ListUserMetaNamespace(ctx context.Context, request *idm.ListUserMetaNamespaceRequest, stream idm.UserMetaService_ListUserMetaNamespaceStream) error {

	defer stream.Close()
	dao := servicecontext.GetDAO(ctx).(meta.DAO).GetNamespaceDao()
	if results, err := dao.List(); err == nil {
		for _, result := range results {
			stream.Send(&idm.ListUserMetaNamespaceResponse{UserMetaNamespace: result})
		}
	}
	return nil
}

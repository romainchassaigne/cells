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

package registry

// ListServicesWithMicroMeta lists the services that have a specific meta name and value associated
func (c *pydioregistry) ListServicesWithMicroMeta(metaName string, metaValue ...string) ([]Service, error) {

	var result []Service

	c.runningmutex.Lock()
	defer c.runningmutex.Unlock()

	for _, rs := range c.running {
		if len(rs.Nodes) == 0 {
			continue
		}
		if value, ok := rs.Nodes[0].Metadata[metaName]; ok {
			// Compare meta value if passed
			if len(metaValue) > 0 && value != metaValue[0] {
				continue
			}
			if rs.Metadata == nil {
				rs.Metadata = make(map[string]string, 1)
			}
			rs.Metadata[metaName] = value

			if service, ok := c.register[rs.Name]; ok {
				result = append(result, service)
			} else {
				result = append(result, &mockService{rs.Name, true, rs.Nodes, []string{}, false})
			}
		}
	}

	return result, nil
}

// ListServicesWithMicroMeta lists the services that have a specific meta name and value associated
func ListServicesWithMicroMeta(metaName string, metaValue ...string) ([]Service, error) {
	return Default.ListServicesWithMicroMeta(metaName, metaValue...)
}

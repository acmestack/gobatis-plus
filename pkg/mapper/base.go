/*
 * Copyright (c) 2022, AcmeStack
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mapper

import "context"

type BaseMapper[T any] interface {
	Insert(ctx context.Context, entity T) int64

	InsertBatch(ctx context.Context, entities ...T) (int64, int64)

	DeleteById(ctx context.Context, id any) int64

	DeleteBatchIds(ctx context.Context, ids []any) int64

	UpdateById(ctx context.Context, entity T) int64

	SelectById(ctx context.Context, id any) T

	SelectBatchIds(ctx context.Context, ids []any) []T

	SelectOne(ctx context.Context, entity T) T

	SelectCount(ctx context.Context, entity T) int64
}

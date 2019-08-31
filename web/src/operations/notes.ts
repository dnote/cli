/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

// This module provides interfaces to perform operations. It abstarcts
// the backend implementation and thus unifies the API for web and desktop clients.

import * as notesService from '../services/notes';
import { NoteData } from './types';
import { Filters } from '../libs/filters';

export function fetch(params: Filters) {
  return notesService.fetch(params);
}

export interface FetchOneParams {
  q?: string;
}

export function fetchOne(noteUUID: string, params: FetchOneParams = {}) {
  return notesService.fetchOne(noteUUID, params);
}

interface CreateParams {
  bookUUID: string;
  content: string;
}

export function create({ bookUUID, content }: CreateParams) {
  return notesService.create({ book_uuid: bookUUID, content });
}

interface UpdateParams {
  book_uuid?: string;
  content?: string;
  public?: boolean;
}

export function update(
  noteUUID: string,
  input: UpdateParams
): Promise<NoteData> {
  return notesService.update(noteUUID, input).then(res => {
    return res.result;
  });
}

export function remove(noteUUID) {
  return notesService.remove(noteUUID);
}
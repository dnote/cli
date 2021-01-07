/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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

package operations

import (
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/models"
	"github.com/pkg/errors"
)

func TestGetNote(t *testing.T) {
	user := models.SetUpUserData()
	anotherUser := models.SetUpUserData()

	defer models.ClearTestData(models.TestDB)

	b1 := models.Book{
		UserID: user.ID,
		Label:  "js",
	}
	models.MustExec(t, models.TestDB.Save(&b1), "preparing b1")

	privateNote := models.Note{
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "privateNote content",
		Deleted:  false,
		Public:   false,
	}
	models.MustExec(t, models.TestDB.Save(&privateNote), "preparing privateNote")

	publicNote := models.Note{
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "privateNote content",
		Deleted:  false,
		Public:   true,
	}
	models.MustExec(t, models.TestDB.Save(&publicNote), "preparing privateNote")

	var privateNoteRecord, publicNoteRecord models.Note
	models.MustExec(t, models.TestDB.Where("uuid = ?", privateNote.UUID).Preload("Book").Preload("User").First(&privateNoteRecord), "finding privateNote")
	models.MustExec(t, models.TestDB.Where("uuid = ?", publicNote.UUID).Preload("Book").Preload("User").First(&publicNoteRecord), "finding publicNote")

	testCases := []struct {
		name         string
		user         models.User
		note         models.Note
		expectedOK   bool
		expectedNote models.Note
	}{
		{
			name:         "owner accessing private note",
			user:         user,
			note:         privateNote,
			expectedOK:   true,
			expectedNote: privateNoteRecord,
		},
		{
			name:         "non-owner accessing private note",
			user:         anotherUser,
			note:         privateNote,
			expectedOK:   false,
			expectedNote: models.Note{},
		},
		{
			name:         "non-owner accessing public note",
			user:         anotherUser,
			note:         publicNote,
			expectedOK:   true,
			expectedNote: publicNoteRecord,
		},
		{
			name:         "guest accessing private note",
			user:         models.User{},
			note:         privateNote,
			expectedOK:   false,
			expectedNote: models.Note{},
		},
		{
			name:         "guest accessing public note",
			user:         models.User{},
			note:         publicNote,
			expectedOK:   true,
			expectedNote: publicNoteRecord,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			note, ok, err := GetNote(models.TestDB, tc.note.UUID, tc.user)
			if err != nil {
				t.Fatal(errors.Wrap(err, "executing"))
			}

			assert.Equal(t, ok, tc.expectedOK, "ok mismatch")
			assert.DeepEqual(t, note, tc.expectedNote, "note mismatch")
		})
	}
}

func TestGetNote_nonexistent(t *testing.T) {
	user := models.SetUpUserData()

	defer models.ClearTestData(models.TestDB)

	b1 := models.Book{
		UserID: user.ID,
		Label:  "js",
	}
	models.MustExec(t, models.TestDB.Save(&b1), "preparing b1")

	n1UUID := "4fd19336-671e-4ff3-8f22-662b80e22edc"
	n1 := models.Note{
		UUID:     n1UUID,
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "n1 content",
		Deleted:  false,
		Public:   false,
	}
	models.MustExec(t, models.TestDB.Save(&n1), "preparing n1")

	nonexistentUUID := "4fd19336-671e-4ff3-8f22-662b80e22edd"
	note, ok, err := GetNote(models.TestDB, nonexistentUUID, user)
	if err != nil {
		t.Fatal(errors.Wrap(err, "executing"))
	}

	assert.Equal(t, ok, false, "ok mismatch")
	assert.DeepEqual(t, note, models.Note{}, "note mismatch")
}

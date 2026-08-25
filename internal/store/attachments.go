package store

import "inventorychain/internal/model"

func (s *Store) SaveAttachment(attachment model.Attachment) error {
	return s.put(AttachmentsBucket, attachment.ID, attachment)
}

func (s *Store) GetAttachment(id string) (model.Attachment, error) {
	var attachment model.Attachment
	err := s.get(AttachmentsBucket, id, &attachment)
	return attachment, err
}

func (s *Store) ListAttachments(recordID string) ([]model.Attachment, error) {
	attachments := make([]model.Attachment, 0)
	err := s.list(AttachmentsBucket, func(data []byte) error {
		var attachment model.Attachment
		if err := unmarshal(data, &attachment); err != nil {
			return err
		}
		if attachment.RecordID == recordID {
			attachments = append(attachments, attachment)
		}
		return nil
	})
	return attachments, err
}

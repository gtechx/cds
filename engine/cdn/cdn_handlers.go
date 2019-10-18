package cdn

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	cdnauth "github.com/ovh/cds/engine/api/authentication/cdn"
	"github.com/ovh/cds/engine/service"
	"github.com/ovh/cds/sdk"
)

// Status returns sdk.MonitoringStatus, implements interface service.Service
func (s *Service) Status() sdk.MonitoringStatus {
	m := s.CommonMonitoring()

	status := sdk.MonitoringStatusOK

	m.Lines = append(m.Lines, sdk.MonitoringStatusLine{Component: "CDN", Value: status, Status: status})

	return m
}

func (s *Service) statusHandler() service.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var status = http.StatusOK
		return service.WriteJSON(w, s.Status(), status)
	}
}

func (s *Service) getDownloadHandler() service.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		token := vars["token"]

		cdnRequest, err := cdnauth.VerifyToken(s.ParsedAPIPublicKey, token)
		if err != nil {
			return sdk.WrapError(sdk.ErrForbidden, "cannot verify token : %v", err)
		}
		w.Header().Add("Content-Type", "application/octet-stream")

		var file io.ReadCloser
		switch cdnRequest.Type {
		case sdk.CDNArtifactType:
			if cdnRequest.Artifact == nil {
				return fmt.Errorf("cannot download artifact, need artifact description in cdn request token")
			}
			w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", cdnRequest.Artifact.Name))

			var err error
			file, err = s.downloadArtifact(r, *cdnRequest)
			if err != nil {
				return sdk.WrapError(err, "cannot download artifact")
			}
		default:
			return fmt.Errorf("cannot download, unknown type %s", cdnRequest.Type)
		}

		if _, err := io.Copy(w, file); err != nil {
			_ = file.Close()
			return sdk.WrapError(err, "Cannot stream file")
		}

		if err := file.Close(); err != nil {
			return sdk.WrapError(err, "Cannot close file")
		}

		return nil
	}
}

func (s *Service) postUploadHandler() service.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		token := vars["token"]

		cdnRequest, err := cdnauth.VerifyToken(s.ParsedAPIPublicKey, token)
		if err != nil {
			return sdk.WrapError(sdk.ErrForbidden, "cannot verify token")
		}

		// decrypt JWT TOKEN
		// Get payload to check which kind of data it is

		// config --> nodeRunID, nodeJobRunID, step for logs, tag, name, projectKey

		switch cdnRequest.Type {
		case sdk.CDNArtifactType:
			if cdnRequest.Artifact == nil {
				return fmt.Errorf("cannot upload artifact, need artifact description in cdn request token")
			}
			artifact, err := s.storeArtifact(r.Body, *cdnRequest)
			if err != nil {
				return sdk.WrapError(err, "cannot store artifact")
			}
			return service.WriteJSON(w, *artifact, http.StatusOK)
		default:
			return fmt.Errorf("cannot download, unknown type %s", cdnRequest.Type)
		}
	}
}

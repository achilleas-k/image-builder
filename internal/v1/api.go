// Package v1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package v1

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// Defines values for Distributions.
const (
	Centos8  Distributions = "centos-8"
	Centos9  Distributions = "centos-9"
	Fedora35 Distributions = "fedora-35"
	Fedora36 Distributions = "fedora-36"
	Fedora37 Distributions = "fedora-37"
	Fedora38 Distributions = "fedora-38"
	Rhel8    Distributions = "rhel-8"
	Rhel84   Distributions = "rhel-84"
	Rhel85   Distributions = "rhel-85"
	Rhel86   Distributions = "rhel-86"
	Rhel87   Distributions = "rhel-87"
	Rhel9    Distributions = "rhel-9"
	Rhel90   Distributions = "rhel-90"
)

// Defines values for ImageRequestArchitecture.
const (
	X8664 ImageRequestArchitecture = "x86_64"
)

// Defines values for ImageStatusStatus.
const (
	ImageStatusStatusBuilding    ImageStatusStatus = "building"
	ImageStatusStatusFailure     ImageStatusStatus = "failure"
	ImageStatusStatusPending     ImageStatusStatus = "pending"
	ImageStatusStatusRegistering ImageStatusStatus = "registering"
	ImageStatusStatusSuccess     ImageStatusStatus = "success"
	ImageStatusStatusUploading   ImageStatusStatus = "uploading"
)

// Defines values for ImageTypes.
const (
	ImageTypesAmi               ImageTypes = "ami"
	ImageTypesAws               ImageTypes = "aws"
	ImageTypesAzure             ImageTypes = "azure"
	ImageTypesEdgeCommit        ImageTypes = "edge-commit"
	ImageTypesEdgeInstaller     ImageTypes = "edge-installer"
	ImageTypesGcp               ImageTypes = "gcp"
	ImageTypesGuestImage        ImageTypes = "guest-image"
	ImageTypesImageInstaller    ImageTypes = "image-installer"
	ImageTypesRhelEdgeCommit    ImageTypes = "rhel-edge-commit"
	ImageTypesRhelEdgeInstaller ImageTypes = "rhel-edge-installer"
	ImageTypesVhd               ImageTypes = "vhd"
	ImageTypesVsphere           ImageTypes = "vsphere"
)

// Defines values for UploadStatusStatus.
const (
	UploadStatusStatusFailure UploadStatusStatus = "failure"
	UploadStatusStatusPending UploadStatusStatus = "pending"
	UploadStatusStatusRunning UploadStatusStatus = "running"
	UploadStatusStatusSuccess UploadStatusStatus = "success"
)

// Defines values for UploadTypes.
const (
	UploadTypesAws   UploadTypes = "aws"
	UploadTypesAwsS3 UploadTypes = "aws.s3"
	UploadTypesAzure UploadTypes = "azure"
	UploadTypesGcp   UploadTypes = "gcp"
)

// AWSEC2Clone defines model for AWSEC2Clone.
type AWSEC2Clone struct {
	// A region as described in
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html#concepts-regions
	Region string `json:"region"`

	// An array of AWS account IDs as described in
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/console_account-alias.html
	ShareWithAccounts *[]string `json:"share_with_accounts,omitempty"`
}

// AWSS3UploadRequestOptions defines model for AWSS3UploadRequestOptions.
type AWSS3UploadRequestOptions = map[string]interface{}

// AWSS3UploadStatus defines model for AWSS3UploadStatus.
type AWSS3UploadStatus struct {
	Url string `json:"url"`
}

// AWSUploadRequestOptions defines model for AWSUploadRequestOptions.
type AWSUploadRequestOptions struct {
	ShareWithAccounts []string `json:"share_with_accounts"`
}

// AWSUploadStatus defines model for AWSUploadStatus.
type AWSUploadStatus struct {
	Ami    string `json:"ami"`
	Region string `json:"region"`
}

// ArchitectureItem defines model for ArchitectureItem.
type ArchitectureItem struct {
	Arch       string   `json:"arch"`
	ImageTypes []string `json:"image_types"`
}

// Architectures defines model for Architectures.
type Architectures = []ArchitectureItem

// AzureUploadRequestOptions defines model for AzureUploadRequestOptions.
type AzureUploadRequestOptions struct {
	// Name of the created image.
	// Must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.
	// The total length is limited to 60 characters.
	ImageName *string `json:"image_name,omitempty"`

	// Name of the resource group where the image should be uploaded.
	ResourceGroup string `json:"resource_group"`

	// ID of subscription where the image should be uploaded.
	SubscriptionId string `json:"subscription_id"`

	// ID of the tenant where the image should be uploaded. This link explains how
	// to find it in the Azure Portal:
	// https://docs.microsoft.com/en-us/azure/active-directory/fundamentals/active-directory-how-to-find-tenant
	TenantId string `json:"tenant_id"`
}

// AzureUploadStatus defines model for AzureUploadStatus.
type AzureUploadStatus struct {
	ImageName string `json:"image_name"`
}

// CloneRequest defines model for CloneRequest.
type CloneRequest interface{}

// CloneResponse defines model for CloneResponse.
type CloneResponse struct {
	Id string `json:"id"`
}

// ClonesResponse defines model for ClonesResponse.
type ClonesResponse struct {
	Data  []ClonesResponseItem `json:"data"`
	Links struct {
		First string `json:"first"`
		Last  string `json:"last"`
	} `json:"links"`
	Meta struct {
		Count int `json:"count"`
	} `json:"meta"`
}

// ClonesResponseItem defines model for ClonesResponseItem.
type ClonesResponseItem struct {
	CreatedAt string      `json:"created_at"`
	Id        string      `json:"id"`
	Request   interface{} `json:"request"`
}

// ComposeMetadata defines model for ComposeMetadata.
type ComposeMetadata struct {
	// ID (hash) of the built commit
	OstreeCommit *string `json:"ostree_commit,omitempty"`

	// Package list including NEVRA
	Packages *[]PackageMetadata `json:"packages,omitempty"`
}

// ComposeRequest defines model for ComposeRequest.
type ComposeRequest struct {
	Customizations *Customizations `json:"customizations,omitempty"`
	Distribution   Distributions   `json:"distribution"`
	ImageName      *string         `json:"image_name,omitempty"`

	// Array of exactly one image request. Having more image requests in one compose is currently not supported.
	ImageRequests []ImageRequest `json:"image_requests"`
}

// ComposeResponse defines model for ComposeResponse.
type ComposeResponse struct {
	Id string `json:"id"`
}

// ComposeStatus defines model for ComposeStatus.
type ComposeStatus struct {
	ImageStatus ImageStatus `json:"image_status"`
}

// ComposeStatusError defines model for ComposeStatusError.
type ComposeStatusError struct {
	Details *interface{} `json:"details,omitempty"`
	Id      int          `json:"id"`
	Reason  string       `json:"reason"`
}

// ComposesResponse defines model for ComposesResponse.
type ComposesResponse struct {
	Data  []ComposesResponseItem `json:"data"`
	Links struct {
		First string `json:"first"`
		Last  string `json:"last"`
	} `json:"links"`
	Meta struct {
		Count int `json:"count"`
	} `json:"meta"`
}

// ComposesResponseItem defines model for ComposesResponseItem.
type ComposesResponseItem struct {
	CreatedAt string      `json:"created_at"`
	Id        string      `json:"id"`
	ImageName *string     `json:"image_name,omitempty"`
	Request   interface{} `json:"request"`
}

// Customizations defines model for Customizations.
type Customizations struct {
	Filesystem          *[]Filesystem `json:"filesystem,omitempty"`
	Packages            *[]string     `json:"packages,omitempty"`
	PayloadRepositories *[]Repository `json:"payload_repositories,omitempty"`
	Subscription        *Subscription `json:"subscription,omitempty"`

	// list of users that a customer can add, also specifying their respective groups and SSH keys
	Users *[]User `json:"users,omitempty"`
}

// DistributionItem defines model for DistributionItem.
type DistributionItem struct {
	Description string `json:"description"`
	Name        string `json:"name"`
}

// Distributions defines model for Distributions.
type Distributions string

// DistributionsResponse defines model for DistributionsResponse.
type DistributionsResponse = []DistributionItem

// Filesystem defines model for Filesystem.
type Filesystem struct {
	MinSize    uint64 `json:"min_size"`
	Mountpoint string `json:"mountpoint"`
}

// GCPUploadRequestOptions defines model for GCPUploadRequestOptions.
type GCPUploadRequestOptions struct {
	// List of valid Google accounts to share the imported Compute Node image with.
	// Each string must contain a specifier of the account type. Valid formats are:
	//   - 'user:{emailid}': An email address that represents a specific
	//     Google account. For example, 'alice@example.com'.
	//   - 'serviceAccount:{emailid}': An email address that represents a
	//     service account. For example, 'my-other-app@appspot.gserviceaccount.com'.
	//   - 'group:{emailid}': An email address that represents a Google group.
	//     For example, 'admins@example.com'.
	//   - 'domain:{domain}': The G Suite domain (primary) that represents all
	//     the users of that domain. For example, 'google.com' or 'example.com'.
	//     If not specified, the imported Compute Node image is not shared with any
	//     account.
	ShareWithAccounts []string `json:"share_with_accounts"`
}

// GCPUploadStatus defines model for GCPUploadStatus.
type GCPUploadStatus struct {
	ImageName string `json:"image_name"`
	ProjectId string `json:"project_id"`
}

// HTTPError defines model for HTTPError.
type HTTPError struct {
	Detail string `json:"detail"`
	Title  string `json:"title"`
}

// HTTPErrorList defines model for HTTPErrorList.
type HTTPErrorList struct {
	Errors []HTTPError `json:"errors"`
}

// ImageRequest defines model for ImageRequest.
type ImageRequest struct {
	// CPU architecture of the image, only x86_64 is currently supported.
	Architecture  ImageRequestArchitecture `json:"architecture"`
	ImageType     ImageTypes               `json:"image_type"`
	Ostree        *OSTree                  `json:"ostree,omitempty"`
	UploadRequest UploadRequest            `json:"upload_request"`
}

// CPU architecture of the image, only x86_64 is currently supported.
type ImageRequestArchitecture string

// ImageStatus defines model for ImageStatus.
type ImageStatus struct {
	Error        *ComposeStatusError `json:"error,omitempty"`
	Status       ImageStatusStatus   `json:"status"`
	UploadStatus *UploadStatus       `json:"upload_status,omitempty"`
}

// ImageStatusStatus defines model for ImageStatus.Status.
type ImageStatusStatus string

// ImageTypes defines model for ImageTypes.
type ImageTypes string

// OSTree defines model for OSTree.
type OSTree struct {
	// A URL which, if set, is used for fetching content. Implies that `url` is set as well,
	// which will be used for metadata only.
	Contenturl *string `json:"contenturl,omitempty"`

	// Can be either a commit (example: 02604b2da6e954bd34b8b82a835e5a77d2b60ffa), or a branch-like reference (example: rhel/8/x86_64/edge)
	Parent *string `json:"parent,omitempty"`
	Ref    *string `json:"ref,omitempty"`

	// Determines whether a valid subscription manager (candlepin) identity is required to
	// access this repository. Consumer certificates will be used as client certificates when
	// fetching metadata and content.
	Rhsm *bool   `json:"rhsm,omitempty"`
	Url  *string `json:"url,omitempty"`
}

// Package defines model for Package.
type Package struct {
	Name    string `json:"name"`
	Summary string `json:"summary"`
}

// PackageMetadata defines model for PackageMetadata.
type PackageMetadata struct {
	Arch      string  `json:"arch"`
	Epoch     *string `json:"epoch,omitempty"`
	Name      string  `json:"name"`
	Release   string  `json:"release"`
	Sigmd5    string  `json:"sigmd5"`
	Signature *string `json:"signature,omitempty"`
	Type      string  `json:"type"`
	Version   string  `json:"version"`
}

// PackagesResponse defines model for PackagesResponse.
type PackagesResponse struct {
	Data  []Package `json:"data"`
	Links struct {
		First string `json:"first"`
		Last  string `json:"last"`
	} `json:"links"`
	Meta struct {
		Count int `json:"count"`
	} `json:"meta"`
}

// Readiness defines model for Readiness.
type Readiness struct {
	Readiness string `json:"readiness"`
}

// Repository defines model for Repository.
type Repository struct {
	Baseurl    *string `json:"baseurl,omitempty"`
	CheckGpg   *bool   `json:"check_gpg,omitempty"`
	Gpgkey     *string `json:"gpgkey,omitempty"`
	IgnoreSsl  *bool   `json:"ignore_ssl,omitempty"`
	Metalink   *string `json:"metalink,omitempty"`
	Mirrorlist *string `json:"mirrorlist,omitempty"`
	Rhsm       bool    `json:"rhsm"`
}

// Subscription defines model for Subscription.
type Subscription struct {
	ActivationKey string `json:"activation-key"`
	BaseUrl       string `json:"base-url"`
	Insights      bool   `json:"insights"`
	Organization  int    `json:"organization"`
	ServerUrl     string `json:"server-url"`
}

// UploadRequest defines model for UploadRequest.
type UploadRequest struct {
	Options interface{} `json:"options"`
	Type    UploadTypes `json:"type"`
}

// UploadStatus defines model for UploadStatus.
type UploadStatus struct {
	Options interface{}        `json:"options"`
	Status  UploadStatusStatus `json:"status"`
	Type    UploadTypes        `json:"type"`
}

// UploadStatusStatus defines model for UploadStatus.Status.
type UploadStatusStatus string

// UploadTypes defines model for UploadTypes.
type UploadTypes string

// User defines model for User.
type User struct {
	Name   string `json:"name"`
	SshKey string `json:"ssh_key"`
}

// Version defines model for Version.
type Version struct {
	Version string `json:"version"`
}

// ComposeImageJSONBody defines parameters for ComposeImage.
type ComposeImageJSONBody = ComposeRequest

// GetComposesParams defines parameters for GetComposes.
type GetComposesParams struct {
	// max amount of composes, default 100
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`

	// composes page offset, default 0
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`
}

// CloneComposeJSONBody defines parameters for CloneCompose.
type CloneComposeJSONBody = CloneRequest

// GetComposeClonesParams defines parameters for GetComposeClones.
type GetComposeClonesParams struct {
	// max amount of composes, default 100
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`

	// composes page offset, default 0
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`
}

// GetPackagesParams defines parameters for GetPackages.
type GetPackagesParams struct {
	// distribution to look up packages for
	Distribution Distributions `form:"distribution" json:"distribution"`

	// architecture to look up packages for
	Architecture GetPackagesParamsArchitecture `form:"architecture" json:"architecture"`

	// packages to look for
	Search string `form:"search" json:"search"`

	// max amount of packages, default 100
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`

	// packages page offset, default 0
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`
}

// GetPackagesParamsArchitecture defines parameters for GetPackages.
type GetPackagesParamsArchitecture string

// ComposeImageJSONRequestBody defines body for ComposeImage for application/json ContentType.
type ComposeImageJSONRequestBody = ComposeImageJSONBody

// CloneComposeJSONRequestBody defines body for CloneCompose for application/json ContentType.
type CloneComposeJSONRequestBody = CloneComposeJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// get the architectures and their image types available for a given distribution
	// (GET /architectures/{distribution})
	GetArchitectures(ctx echo.Context, distribution string) error
	// get status of a compose clone
	// (GET /clones/{id})
	GetCloneStatus(ctx echo.Context, id string) error
	// compose image
	// (POST /compose)
	ComposeImage(ctx echo.Context) error
	// get a collection of previous compose requests for the logged in user
	// (GET /composes)
	GetComposes(ctx echo.Context, params GetComposesParams) error
	// get status of an image compose
	// (GET /composes/{composeId})
	GetComposeStatus(ctx echo.Context, composeId string) error
	// clone a compose
	// (POST /composes/{composeId}/clone)
	CloneCompose(ctx echo.Context, composeId string) error
	// get clones of a compose
	// (GET /composes/{composeId}/clones)
	GetComposeClones(ctx echo.Context, composeId string, params GetComposeClonesParams) error
	// get metadata of an image compose
	// (GET /composes/{composeId}/metadata)
	GetComposeMetadata(ctx echo.Context, composeId string) error
	// get the available distributions
	// (GET /distributions)
	GetDistributions(ctx echo.Context) error
	// get the openapi json specification
	// (GET /openapi.json)
	GetOpenapiJson(ctx echo.Context) error

	// (GET /packages)
	GetPackages(ctx echo.Context, params GetPackagesParams) error
	// return the readiness
	// (GET /ready)
	GetReadiness(ctx echo.Context) error
	// get the service version
	// (GET /version)
	GetVersion(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetArchitectures converts echo context to params.
func (w *ServerInterfaceWrapper) GetArchitectures(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "distribution" -------------
	var distribution string

	err = runtime.BindStyledParameterWithLocation("simple", false, "distribution", runtime.ParamLocationPath, ctx.Param("distribution"), &distribution)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter distribution: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetArchitectures(ctx, distribution)
	return err
}

// GetCloneStatus converts echo context to params.
func (w *ServerInterfaceWrapper) GetCloneStatus(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetCloneStatus(ctx, id)
	return err
}

// ComposeImage converts echo context to params.
func (w *ServerInterfaceWrapper) ComposeImage(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ComposeImage(ctx)
	return err
}

// GetComposes converts echo context to params.
func (w *ServerInterfaceWrapper) GetComposes(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetComposesParams
	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetComposes(ctx, params)
	return err
}

// GetComposeStatus converts echo context to params.
func (w *ServerInterfaceWrapper) GetComposeStatus(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "composeId" -------------
	var composeId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "composeId", runtime.ParamLocationPath, ctx.Param("composeId"), &composeId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter composeId: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetComposeStatus(ctx, composeId)
	return err
}

// CloneCompose converts echo context to params.
func (w *ServerInterfaceWrapper) CloneCompose(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "composeId" -------------
	var composeId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "composeId", runtime.ParamLocationPath, ctx.Param("composeId"), &composeId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter composeId: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CloneCompose(ctx, composeId)
	return err
}

// GetComposeClones converts echo context to params.
func (w *ServerInterfaceWrapper) GetComposeClones(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "composeId" -------------
	var composeId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "composeId", runtime.ParamLocationPath, ctx.Param("composeId"), &composeId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter composeId: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetComposeClonesParams
	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetComposeClones(ctx, composeId, params)
	return err
}

// GetComposeMetadata converts echo context to params.
func (w *ServerInterfaceWrapper) GetComposeMetadata(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "composeId" -------------
	var composeId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "composeId", runtime.ParamLocationPath, ctx.Param("composeId"), &composeId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter composeId: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetComposeMetadata(ctx, composeId)
	return err
}

// GetDistributions converts echo context to params.
func (w *ServerInterfaceWrapper) GetDistributions(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetDistributions(ctx)
	return err
}

// GetOpenapiJson converts echo context to params.
func (w *ServerInterfaceWrapper) GetOpenapiJson(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetOpenapiJson(ctx)
	return err
}

// GetPackages converts echo context to params.
func (w *ServerInterfaceWrapper) GetPackages(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPackagesParams
	// ------------- Required query parameter "distribution" -------------

	err = runtime.BindQueryParameter("form", true, true, "distribution", ctx.QueryParams(), &params.Distribution)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter distribution: %s", err))
	}

	// ------------- Required query parameter "architecture" -------------

	err = runtime.BindQueryParameter("form", true, true, "architecture", ctx.QueryParams(), &params.Architecture)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter architecture: %s", err))
	}

	// ------------- Required query parameter "search" -------------

	err = runtime.BindQueryParameter("form", true, true, "search", ctx.QueryParams(), &params.Search)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter search: %s", err))
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPackages(ctx, params)
	return err
}

// GetReadiness converts echo context to params.
func (w *ServerInterfaceWrapper) GetReadiness(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetReadiness(ctx)
	return err
}

// GetVersion converts echo context to params.
func (w *ServerInterfaceWrapper) GetVersion(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetVersion(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/architectures/:distribution", wrapper.GetArchitectures)
	router.GET(baseURL+"/clones/:id", wrapper.GetCloneStatus)
	router.POST(baseURL+"/compose", wrapper.ComposeImage)
	router.GET(baseURL+"/composes", wrapper.GetComposes)
	router.GET(baseURL+"/composes/:composeId", wrapper.GetComposeStatus)
	router.POST(baseURL+"/composes/:composeId/clone", wrapper.CloneCompose)
	router.GET(baseURL+"/composes/:composeId/clones", wrapper.GetComposeClones)
	router.GET(baseURL+"/composes/:composeId/metadata", wrapper.GetComposeMetadata)
	router.GET(baseURL+"/distributions", wrapper.GetDistributions)
	router.GET(baseURL+"/openapi.json", wrapper.GetOpenapiJson)
	router.GET(baseURL+"/packages", wrapper.GetPackages)
	router.GET(baseURL+"/ready", wrapper.GetReadiness)
	router.GET(baseURL+"/version", wrapper.GetVersion)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}

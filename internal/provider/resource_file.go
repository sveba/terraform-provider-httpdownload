package provider

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/svetob/terraform-provider-httpdownload/internal/client"
)

var (
	_ resource.Resource              = &httpDownloadFileResource{}
	_ resource.ResourceWithConfigure = &httpDownloadFileResource{}
)

func NewHttpDownloadFileResource() resource.Resource {
	return &httpDownloadFileResource{}
}

type httpDownloadFileResource struct {
	client *client.HttpDownloadClient
}

func (r *httpDownloadFileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (r *httpDownloadFileResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Represents a download of a file from a HTTP server.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Required:    true,
				Description: "Url of the file to download.",
			},
			"dest": schema.StringAttribute{
				Required:    true,
				Description: "Destination file path.",
			},
			"checksum": schema.StringAttribute{
				Computed:    true,
				Description: "Destination file path.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *httpDownloadFileResource) Configure(ctx context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.HttpDownloadClient)
}

func (r httpDownloadFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan HttpDownloadFile
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DownloadFileToDest(ctx, plan.Dest.ValueString(), plan.Url.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Resource",
			err.Error(),
		)
		return
	}

	checksum_on_disk, err := getSha256Hash(plan.Dest.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Resource",
			err.Error(),
		)
		return
	}

	plan.Checksum = types.StringValue(checksum_on_disk)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r httpDownloadFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state HttpDownloadFile
	req.State.Get(ctx, &state)

	checksum_on_disk, err := getSha256Hash(state.Dest.ValueString())
	if err != nil || checksum_on_disk != state.Checksum.ValueString() {
		resp.State.RemoveResource(ctx)
		return
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r httpDownloadFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r httpDownloadFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state HttpDownloadFile
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := os.Remove(state.Dest.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Resource",
			err.Error(),
		)
		return
	}
}

func getSha256Hash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

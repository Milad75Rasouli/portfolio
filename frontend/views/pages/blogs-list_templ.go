// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.663
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "github.com/Milad75Rasouli/portfolio/frontend/views/layouts"
import "github.com/Milad75Rasouli/portfolio/internal/model"
import "fmt"

func List(blogs []model.Blog) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!-- return list -->")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func BlogList(blogs string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = layouts.Base("Blog List", 2, fmt.Sprintf(`
  <div class="b-content-divider b-content-divider-day b-content-vr col p-0">
    <div class="container h-100">
        <div class="h-100 align-hv-center">
            <div class="pt-3 pb-3 p-4 border border-light bg-opacity-75 text-start bg-dark blogs-list overflow-auto">
              <div class="row row-cols-1 row-cols-xxl-3 row-cols-xl-3 row-cols-md-2 row-cols-lg-2 g-4">                
                  %s
<!-- 
                <div class="col">
                  <div class="card bg-light text-dark border-light">
                    <img src="/static/img/blog.png" class="card-img-top" alt="...">
                    <div class="card-body">
                      <h5 class="card-title">Card title</h5>
                      <p class="card-text">This is a longer card with supporting text below as a natural lead-in to additional content. This content is a little bit longer.</p>
                      <a href="#" class="btn btn-dark float-start">Read</a>
                    </div>
                    <div class="card-footer text-dark border-light">
                      <div class="float-start text-small">
                        2 days ago 
                      </div> 
                      <div class="badge p-2 text-bg-primary text-wrap float-end" style="width: 6rem;">
                        database
                      </div>
                    </div>
                  </div>
                </div> -->

              </div>
            </div>
        </div>
        </div>
    </div>
  </div>
  `, blogs)).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

// Code generated by qtc from "jeniskeluhan.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line jeniskeluhan.qtpl:1
package view

//line jeniskeluhan.qtpl:1
import (
	"b2t_helpdesk/injector"
	"encoding/json"
	"github.com/valyala/fasthttp"
)

//line jeniskeluhan.qtpl:7
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line jeniskeluhan.qtpl:7
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line jeniskeluhan.qtpl:8
type JenisKeluhan struct {
	CTX       *fasthttp.RequestCtx
	Dinjector *injector.Injector
}

//line jeniskeluhan.qtpl:15
func (p *JenisKeluhan) StreamTitle(qw422016 *qt422016.Writer) {
//line jeniskeluhan.qtpl:15
	qw422016.N().S(`
	Daftar Jenis Keluhan
`)
//line jeniskeluhan.qtpl:17
}

//line jeniskeluhan.qtpl:17
func (p *JenisKeluhan) WriteTitle(qq422016 qtio422016.Writer) {
//line jeniskeluhan.qtpl:17
	qw422016 := qt422016.AcquireWriter(qq422016)
//line jeniskeluhan.qtpl:17
	p.StreamTitle(qw422016)
//line jeniskeluhan.qtpl:17
	qt422016.ReleaseWriter(qw422016)
//line jeniskeluhan.qtpl:17
}

//line jeniskeluhan.qtpl:17
func (p *JenisKeluhan) Title() string {
//line jeniskeluhan.qtpl:17
	qb422016 := qt422016.AcquireByteBuffer()
//line jeniskeluhan.qtpl:17
	p.WriteTitle(qb422016)
//line jeniskeluhan.qtpl:17
	qs422016 := string(qb422016.B)
//line jeniskeluhan.qtpl:17
	qt422016.ReleaseByteBuffer(qb422016)
//line jeniskeluhan.qtpl:17
	return qs422016
//line jeniskeluhan.qtpl:17
}

//line jeniskeluhan.qtpl:20
func (p *JenisKeluhan) StreamBody(qw422016 *qt422016.Writer) {
//line jeniskeluhan.qtpl:20
	qw422016.N().S(`
	<div class="card shadow mb-4">
		<div class="card-body">
			<button class="btn btn-success" type="button" onclick="readytambah()" data-toggle="modal" data-target="#addModal">Tambah Jenis</button>
			<div class="table-responsive">
				<table class="table table-bordered" width="100%" cellspacing="0">
					<thead>
						<tr>
							<th>#</th>
							<th>ID</th>
							<th>Jenis keluhan</th>
							<th>Hint</th>
							<th>Lanjut Ke Bantuan Langsung</th>
						</tr>
					</thead>
					<tbody>
					`)
//line jeniskeluhan.qtpl:37
	results, _ := p.Dinjector.DB.Query("SELECT id, jenis, hint, autoinput FROM jenisticket ORDER BY id ASC")
	type keluhan struct {
		Id        int    `json:"id"`
		Jenis     string `json:"jenis"`
		Hint      string `json:"hint"`
		Autoinput int    `json:"autoinput"`
	}
	var daftarkeluhan []keluhan
	for results.Next() {
		var buffkeluhan keluhan
		err := results.Scan(&buffkeluhan.Id, &buffkeluhan.Jenis, &buffkeluhan.Hint, &buffkeluhan.Autoinput)
		if err == nil {
			daftarkeluhan = append(daftarkeluhan, buffkeluhan)
		}
	}

//line jeniskeluhan.qtpl:52
	qw422016.N().S(`
					
					`)
//line jeniskeluhan.qtpl:54
	for _, value := range daftarkeluhan {
//line jeniskeluhan.qtpl:54
		qw422016.N().S(`
						<tr>
							`)
//line jeniskeluhan.qtpl:57
		dtjson, _ := json.Marshal(value)

//line jeniskeluhan.qtpl:58
		qw422016.N().S(`
							<td>
								<button class="btn btn-warning" data-toggle="modal" data-target="#addModal" onclick="readyEdit('`)
//line jeniskeluhan.qtpl:60
		qw422016.E().S(string(dtjson))
//line jeniskeluhan.qtpl:60
		qw422016.N().S(`')" type="button"><i class="fa fa-edit"></i></button>
								<a href="/jeniskeluhan/delete/`)
//line jeniskeluhan.qtpl:61
		qw422016.N().D(value.Id)
//line jeniskeluhan.qtpl:61
		qw422016.N().S(`"><button class="btn btn-danger" type="button"><i class="fa fa-trash"></i></button></a>
							</td>
							<td>`)
//line jeniskeluhan.qtpl:63
		qw422016.N().D(value.Id)
//line jeniskeluhan.qtpl:63
		qw422016.N().S(`</td>
							<td>`)
//line jeniskeluhan.qtpl:64
		qw422016.E().S(value.Jenis)
//line jeniskeluhan.qtpl:64
		qw422016.N().S(`</td>
							<td>`)
//line jeniskeluhan.qtpl:65
		qw422016.E().S(value.Hint)
//line jeniskeluhan.qtpl:65
		qw422016.N().S(`</td>
							<td>`)
//line jeniskeluhan.qtpl:66
		qw422016.N().D(value.Autoinput)
//line jeniskeluhan.qtpl:66
		qw422016.N().S(`</td>
						</tr>
					`)
//line jeniskeluhan.qtpl:68
	}
//line jeniskeluhan.qtpl:68
	qw422016.N().S(`
					</tbody>
				</table>
			</div>
		</div>
	</div>

`)
//line jeniskeluhan.qtpl:75
}

//line jeniskeluhan.qtpl:75
func (p *JenisKeluhan) WriteBody(qq422016 qtio422016.Writer) {
//line jeniskeluhan.qtpl:75
	qw422016 := qt422016.AcquireWriter(qq422016)
//line jeniskeluhan.qtpl:75
	p.StreamBody(qw422016)
//line jeniskeluhan.qtpl:75
	qt422016.ReleaseWriter(qw422016)
//line jeniskeluhan.qtpl:75
}

//line jeniskeluhan.qtpl:75
func (p *JenisKeluhan) Body() string {
//line jeniskeluhan.qtpl:75
	qb422016 := qt422016.AcquireByteBuffer()
//line jeniskeluhan.qtpl:75
	p.WriteBody(qb422016)
//line jeniskeluhan.qtpl:75
	qs422016 := string(qb422016.B)
//line jeniskeluhan.qtpl:75
	qt422016.ReleaseByteBuffer(qb422016)
//line jeniskeluhan.qtpl:75
	return qs422016
//line jeniskeluhan.qtpl:75
}

//line jeniskeluhan.qtpl:77
func (p *JenisKeluhan) StreamModal(qw422016 *qt422016.Writer) {
//line jeniskeluhan.qtpl:77
	qw422016.N().S(`
	<div class="modal fade" id="addModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalLabel" aria-hidden="true">
        <div class="modal-dialog modal-xl" role="document">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title" id="exampleModalLabel">Tambah / Edit Jenis Keluhan</span></h5>
                    <button class="close" type="button" data-dismiss="modal" aria-label="Close">
                        <span aria-hidden="true">×</span>
                    </button>
                </div>
				<form action="/jeniskeluhan/update" method="post">
					<div class="modal-body">
						<div class="row">
							<div class="col-12">
								<div class="form group">
									<input name="id" id="id" type="hidden">
									<label>Jenis</label>
									<input name="jenis" id="jenis" class="form-control" type="text">
									<label>Hint</label>
									<textarea name="hint" id="hint" class="form-control" rows="3"></textarea>
									<label>Lanjut Ke Bantuan Langsung (0 Tidak, 1 Iya)</label>
									<input name="autoinput" id="autoinput" class="form-control" type="number" max="1" min="0">
								</div>
							</div>
						</div>
					</div>
					<div class="modal-footer">
						<button class="btn btn-success" type="submit">Simpan</button>
					</div>
                </form>
            </div>
        </div>
    </div>
`)
//line jeniskeluhan.qtpl:110
}

//line jeniskeluhan.qtpl:110
func (p *JenisKeluhan) WriteModal(qq422016 qtio422016.Writer) {
//line jeniskeluhan.qtpl:110
	qw422016 := qt422016.AcquireWriter(qq422016)
//line jeniskeluhan.qtpl:110
	p.StreamModal(qw422016)
//line jeniskeluhan.qtpl:110
	qt422016.ReleaseWriter(qw422016)
//line jeniskeluhan.qtpl:110
}

//line jeniskeluhan.qtpl:110
func (p *JenisKeluhan) Modal() string {
//line jeniskeluhan.qtpl:110
	qb422016 := qt422016.AcquireByteBuffer()
//line jeniskeluhan.qtpl:110
	p.WriteModal(qb422016)
//line jeniskeluhan.qtpl:110
	qs422016 := string(qb422016.B)
//line jeniskeluhan.qtpl:110
	qt422016.ReleaseByteBuffer(qb422016)
//line jeniskeluhan.qtpl:110
	return qs422016
//line jeniskeluhan.qtpl:110
}

//line jeniskeluhan.qtpl:112
func (p *JenisKeluhan) StreamScript(qw422016 *qt422016.Writer) {
//line jeniskeluhan.qtpl:112
	qw422016.N().S(`
<script>
	function readytambah(){
		$('#id').val("0");
		$('#jenis').val("");
		$('#hint').val("");
		$('#autoinput').val("0");
	}

	function readyEdit(d){

		data = JSON.parse(d);

		$('#id').val(data.id);
		$('#jenis').val(data.jenis);
		$('#hint').val(data.hint);
		$('#autoinput').val(data.autoinput);

	}
</script>
`)
//line jeniskeluhan.qtpl:132
}

//line jeniskeluhan.qtpl:132
func (p *JenisKeluhan) WriteScript(qq422016 qtio422016.Writer) {
//line jeniskeluhan.qtpl:132
	qw422016 := qt422016.AcquireWriter(qq422016)
//line jeniskeluhan.qtpl:132
	p.StreamScript(qw422016)
//line jeniskeluhan.qtpl:132
	qt422016.ReleaseWriter(qw422016)
//line jeniskeluhan.qtpl:132
}

//line jeniskeluhan.qtpl:132
func (p *JenisKeluhan) Script() string {
//line jeniskeluhan.qtpl:132
	qb422016 := qt422016.AcquireByteBuffer()
//line jeniskeluhan.qtpl:132
	p.WriteScript(qb422016)
//line jeniskeluhan.qtpl:132
	qs422016 := string(qb422016.B)
//line jeniskeluhan.qtpl:132
	qt422016.ReleaseByteBuffer(qb422016)
//line jeniskeluhan.qtpl:132
	return qs422016
//line jeniskeluhan.qtpl:132
}
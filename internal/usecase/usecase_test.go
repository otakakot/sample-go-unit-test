package usecase_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	"github.com/otakakot/sample-go-unit-test/internal/errors"
	"github.com/otakakot/sample-go-unit-test/internal/model"
	"github.com/otakakot/sample-go-unit-test/internal/repository"
	"github.com/otakakot/sample-go-unit-test/internal/usecase"
)

func TestUsecaseCreate(t *testing.T) {
	t.Parallel()

	type fields struct {
		repository func(*testing.T) repository.Repository
	}

	type args struct {
		ctx  context.Context
		name string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Model
		wantErr bool
	}{
		{
			name: "作成に成功", // テストケースは日本語で書く
			fields: fields{
				repository: func(t *testing.T) repository.Repository {
					t.Helper()                      // おまじない
					ctrl := gomock.NewController(t) // おまじない
					mock := repository.NewMockRepository(ctrl)
					mock.EXPECT().Save(
						gomock.Any(), // <- contextなので任意の値にしておく
						gomock.Any(), // <- model.Model{}はusecase内で生成されるので一旦任意の値にしておく(比較対象の値をごにょごにょしたい)
					).Return(
						nil, // <- ここでSaveメソッドの返り値を指定する error は発生しないので nil を返す
					).Do(func(ctx context.Context, got model.Model) { // <- Do()を用いて Save()に与えられる model.Model{} を検証する
						if _, err := uuid.Parse(got.ID); err != nil { // <- IDが uuid で採番されているか検証する
							t.Errorf("failed to parse uuid")
						}
						opts := []cmp.Option{
							cmpopts.IgnoreFields(model.Model{}, "ID"), // <- IDを比較対象から除外する
						}
						want := model.Model{
							// IDは比較対象から除外しているので指定しない
							Name: "otakakot",
						}
						if diff := cmp.Diff(got, want, opts...); diff != "" {
							t.Errorf("Save() = %v, want %v", got, want)
						}
					})
					// mock.EXPECT().Find()は定義しない = Create() では呼ばれないのだとわかる
					return mock
				},
			},
			args: args{
				ctx:  context.Background(),
				name: "otakakot",
			},
			want: model.Model{
				// ID は usecase 内で採番されるのでここでは指定しない
				Name: "otakakot",
			},
			wantErr: false,
		},
		{
			name: "作成に失敗",
			fields: fields{
				repository: func(t *testing.T) repository.Repository {
					t.Helper()                      // おまじない
					ctrl := gomock.NewController(t) // おまじない
					mock := repository.NewMockRepository(ctrl)
					mock.EXPECT().Save(
						gomock.Any(), // <- contextなので任意の値にしておく
						gomock.Any(), // <- usecase内で生成されるので一旦任意に値にしておく
					).Return(
						fmt.Errorf("failed to save"), // <- Save()メソッドでエラーが発生したと仮定する 独自で定義した errors パッケージを利用したいので fmt.Errorf() を使用
					).Do(func(ctx context.Context, got model.Model) { // <- error は発生するが Save() メソッドに適切な値が渡っているかは検証する
						if _, err := uuid.Parse(got.ID); err != nil {
							t.Errorf("failed to parse uuid")
						}
						opts := []cmp.Option{
							cmpopts.IgnoreFields(model.Model{}, "ID"),
						}
						want := model.Model{
							Name: "otakakot",
						}
						if diff := cmp.Diff(got, want, opts...); diff != "" {
							t.Errorf("Save() = %v, want %v", got, want)
						}
					})
					return mock
				},
			},
			args: args{
				ctx:  context.Background(),
				name: "otakakot",
			},
			want:    model.Model{}, // 空structが返ってくる。ここはポインタ返しておくようにしておく方が interface の設計として適切なんじゃないかと最近思っている
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt // <- おまじない
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			uc := usecase.New(tt.fields.repository(t))
			got, err := uc.Create(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got.ID = "" // <- 比較で落ちないように空文字いれておく
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Usecase.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecaseRead(t *testing.T) {
	t.Parallel()

	type fields struct {
		repository func(*testing.T) repository.Repository
	}

	type args struct {
		ctx context.Context
		id  string
	}

	id := uuid.NewString() // 比較に使いたいので外で宣言しておく

	tests := []struct {
		name     string
		fields   fields
		args     args
		want     model.Model
		wantErr  bool
		checkErr func(t *testing.T, err error) // <- model.NotFoundError を検証したい
	}{
		{
			name: "見つかる",
			fields: fields{
				repository: func(*testing.T) repository.Repository {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockRepository(ctrl)
					mock.EXPECT().Find(
						gomock.Any(), // <- context
						id,           // <- 外で定義している id が渡っていくることを想定
					).Return(
						model.Model{ // <- repository.Find() で返す値の定積
							ID:   id,         // 整合性取るために受け取ったidを指定しておく
							Name: "otakakot", // 期待値と同じ値を定義しておく
						},
						nil,
					) // Find()メソッドの引数で比較ができるため .Do() の実装は不要
					// mock.EXPECT().Save()は定義しない = Read() では呼ばれないのだとわかる
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				id:  id,
			},
			want: model.Model{
				ID:   id,
				Name: "otakakot",
			},
			wantErr: false,
			// checkErr: func(t *testing.T, err error) {} <- 呼ばれないので指定しない
		},
		{
			name: "見つからない",
			fields: fields{
				repository: func(*testing.T) repository.Repository {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockRepository(ctrl)
					mock.EXPECT().Find(
						gomock.Any(), // <- context
						id,           // メソッドは失敗する想定だが適切なidが渡ってくるかは検証する
					).Return(
						model.Model{},
						errors.NewNotFoundError(fmt.Errorf("not found")),
					)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				id:  id,
			},
			want:    model.Model{},
			wantErr: true,
			checkErr: func(t *testing.T, err error) {
				t.Helper()                        // <- おまじない
				if !errors.AsNotFoundError(err) { // usecase.Read() から想定通りの error 型が返ってくるか検証する
					t.Errorf("failed to assert NotFoundError")
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			uc := usecase.New(tt.fields.repository(t))
			got, err := uc.Read(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Usecase.Read() = %v, want %v", got, tt.want)
			}
			if !tt.wantErr { // <- error は発生しない想定なのでここで終了する
				return
			}
			tt.checkErr(t, err)
		})
	}
}

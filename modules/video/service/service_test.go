package service

import (
	"context"
	"errors"
	"testing"

	"github.com/NTHU-LSALAB/NTHU-Distributed-System/modules/video/dao"
	"github.com/NTHU-LSALAB/NTHU-Distributed-System/modules/video/mock/daomock"
	"github.com/NTHU-LSALAB/NTHU-Distributed-System/modules/video/pb"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Service")
}

var (
	errMongoUnknown = errors.New("unknown mongo error")
)

var _ = Describe("Service", func() {
	var (
		controller *gomock.Controller
		videoDAO   *daomock.MockVideoDAO
		svc        *service
		ctx        context.Context
	)

	BeforeEach(func() {
		controller = gomock.NewController(GinkgoT())
		videoDAO = daomock.NewMockVideoDAO(controller)
		svc = NewService(videoDAO)
		ctx = context.Background()
	})

	AfterEach(func() {
		controller.Finish()
	})

	Describe("GetVideo", func() {
		var (
			req  *pb.GetVideoRequest
			id   primitive.ObjectID
			resp *pb.GetVideoResponse
			err  error
		)

		BeforeEach(func() {
			id = primitive.NewObjectID()
			req = &pb.GetVideoRequest{Id: id.Hex()}
		})

		JustBeforeEach(func() {
			resp, err = svc.GetVideo(ctx, req)
		})

		When("mongo error", func() {
			BeforeEach(func() {
				videoDAO.EXPECT().Get(ctx, id).Return(nil, errMongoUnknown)
			})

			It("returns the error", func() {
				Expect(resp).To(BeNil())
				Expect(err).To(MatchError(errMongoUnknown))
			})
		})

		When("video not found", func() {
			BeforeEach(func() {
				videoDAO.EXPECT().Get(ctx, id).Return(nil, dao.ErrVideoNotFound)
			})

			It("returns video not found error", func() {
				Expect(resp).To(BeNil())
				Expect(err).To(MatchError(ErrVideoNotFound))
			})
		})

		When("success", func() {
			var video *dao.Video

			BeforeEach(func() {
				video = dao.NewFakeVideo()
				videoDAO.EXPECT().Get(ctx, id).Return(video, nil)
			})

			It("returns the video with no error", func() {
				Expect(resp).To(Equal(&pb.GetVideoResponse{
					Video: video.ToProto(),
				}))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("ListVideo", func() {
		var (
			req  *pb.ListVideoRequest
			resp *pb.ListVideoResponse
			err  error
		)

		JustBeforeEach(func() {
			resp, err = svc.ListVideo(ctx, req)
		})

		When("mongo error", func() {
			BeforeEach(func() {
				videoDAO.EXPECT().List(ctx, req.GetLimit(), req.GetSkip()).Return(nil, errMongoUnknown)
			})

			It("returns the error", func() {
				Expect(resp).To(BeNil())
				Expect(err).To(MatchError(errMongoUnknown))
			})
		})

		When("success", func() {
			var videos []*dao.Video

			BeforeEach(func() {
				videos = []*dao.Video{dao.NewFakeVideo(), dao.NewFakeVideo()}
				videoDAO.EXPECT().List(ctx, req.GetLimit(), req.GetSkip()).Return(videos, nil)
			})

			It("returns videos with no error", func() {
				Expect(resp).To(Equal(&pb.ListVideoResponse{
					Videos: []*pb.VideoInfo{
						videos[0].ToProto(),
						videos[1].ToProto(),
					},
				}))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("DeleteVideo", func() {
		var (
			req  *pb.DeleteVideoRequest
			id   primitive.ObjectID
			resp *pb.DeleteVideoResponse
			err  error
		)

		BeforeEach(func() {
			id = primitive.NewObjectID()
			req = &pb.DeleteVideoRequest{Id: id.Hex()}
		})

		JustBeforeEach(func() {
			resp, err = svc.DeleteVideo(ctx, req)
		})

		When("mongo error", func() {
			BeforeEach(func() {
				videoDAO.EXPECT().Delete(ctx, id).Return(errMongoUnknown)
			})

			It("returns the error", func() {
				Expect(resp).To(BeNil())
				Expect(err).To(MatchError(errMongoUnknown))
			})
		})

		When("video not found", func() {
			BeforeEach(func() {
				videoDAO.EXPECT().Delete(ctx, id).Return(dao.ErrVideoNotFound)
			})

			It("returns video not found error", func() {
				Expect(resp).To(BeNil())
				Expect(err).To(MatchError(ErrVideoNotFound))
			})
		})

		When("success", func() {
			BeforeEach(func() {
				videoDAO.EXPECT().Delete(ctx, id).Return(nil)
			})

			It("returns no error", func() {
				Expect(resp).To(Equal(&pb.DeleteVideoResponse{}))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
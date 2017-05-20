// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package tests

import (
	"testing"

	"github.com/npiganeau/yep/pool"
	"github.com/npiganeau/yep/yep/models"
	"github.com/npiganeau/yep/yep/models/security"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGroupLoading(t *testing.T) {
	models.SimulateInNewEnvironment(security.SuperUserID, func(env models.Environment) {
		var (
			adminUser, user                    pool.ResUsersSet
			adminGrp, someGroup, everyoneGroup pool.ResGroupsSet
		)
		Convey("Testing Group Loading", t, func() {
			pool.ResGroups().NewSet(env).ReloadGroups()
			groups := pool.ResGroups().NewSet(env).FetchAll()
			So(groups.Len(), ShouldEqual, len(security.Registry.AllGroups()))
			adminUser = pool.ResUsers().Search(env, pool.ResUsers().ID().Equals(security.SuperUserID))
			adminGrp = pool.ResGroups().Search(env, pool.ResGroups().GroupID().Equals(security.GroupAdminID))
			everyoneGroup = pool.ResGroups().Search(env, pool.ResGroups().GroupID().Equals(security.GroupEveryoneID))
			So(adminUser.Groups().Ids(), ShouldHaveLength, 2)
			So(adminUser.Groups().Ids(), ShouldContain, adminGrp.ID())
			So(adminUser.Groups().Ids(), ShouldContain, everyoneGroup.ID())
		})
		Convey("Testing Group ReLoading with a new group", t, func() {
			security.Registry.NewGroup("some_group", "Some Group")
			pool.ResGroups().NewSet(env).ReloadGroups()
			groups := pool.ResGroups().NewSet(env).FetchAll()
			So(groups.Len(), ShouldEqual, len(security.Registry.AllGroups()))
		})
		Convey("Creating a new user with a new group", t, func() {
			adminUser = pool.ResUsers().Search(env, pool.ResUsers().ID().Equals(security.SuperUserID))
			adminGrp = pool.ResGroups().Search(env, pool.ResGroups().GroupID().Equals(security.GroupAdminID))
			everyoneGroup = pool.ResGroups().Search(env, pool.ResGroups().GroupID().Equals(security.GroupEveryoneID))
			someGroup = pool.ResGroups().Search(env, pool.ResGroups().GroupID().Equals("some_group"))
			user = pool.ResUsers().Create(env, &pool.ResUsersData{
				Name:   "Test User",
				Login:  "test_user",
				Groups: someGroup,
			})
			So(user.Groups().Ids(), ShouldHaveLength, 1)
			So(user.Groups().Ids(), ShouldContain, someGroup.ID())
			So(adminUser.Groups().Ids(), ShouldHaveLength, 2)
			So(adminUser.Groups().Ids(), ShouldContain, adminGrp.ID())
			So(adminUser.Groups().Ids(), ShouldContain, everyoneGroup.ID())
		})
		Convey("Giving the new user admin rights", t, func() {
			groups := someGroup.Union(adminGrp)
			user.SetGroups(groups)
			pool.ResGroups().NewSet(env).ReloadGroups()
			adminGrp = pool.ResGroups().Search(env, pool.ResGroups().GroupID().Equals(security.GroupAdminID))
			everyoneGroup = pool.ResGroups().Search(env, pool.ResGroups().GroupID().Equals(security.GroupEveryoneID))
			someGroup = pool.ResGroups().Search(env, pool.ResGroups().GroupID().Equals("some_group"))
			So(user.Groups().Ids(), ShouldHaveLength, 3)
			So(user.Groups().Ids(), ShouldContain, someGroup.ID())
			So(user.Groups().Ids(), ShouldContain, adminGrp.ID())
			So(user.Groups().Ids(), ShouldContain, everyoneGroup.ID())
			So(adminUser.Groups().Ids(), ShouldHaveLength, 2)
			So(adminUser.Groups().Ids(), ShouldContain, adminGrp.ID())
			So(adminUser.Groups().Ids(), ShouldContain, everyoneGroup.ID())
		})
		Convey("Removing rights and checking that after reload, we get admin right", t, func() {
			user = pool.ResUsers().Search(env, pool.ResUsers().Login().Equals("test_user"))
			user.SetGroups(pool.ResGroups().NewSet(env))
			adminUser = pool.ResUsers().Search(env, pool.ResUsers().ID().Equals(security.SuperUserID))
			adminUser.SetGroups(pool.ResGroups().NewSet(env))
			So(user.Groups().Ids(), ShouldBeEmpty)
			So(adminUser.Groups().Ids(), ShouldBeEmpty)
			pool.ResGroups().NewSet(env).ReloadGroups()
			// We need to reload admin and everyone groups because they changed id
			adminGrp = pool.ResGroups().Search(env, pool.ResGroups().GroupID().Equals(security.GroupAdminID))
			everyoneGroup = pool.ResGroups().Search(env, pool.ResGroups().GroupID().Equals(security.GroupEveryoneID))
			So(user.Groups().Ids(), ShouldHaveLength, 1)
			So(user.Groups().Ids(), ShouldContain, everyoneGroup.ID())
			So(adminUser.Groups().Ids(), ShouldHaveLength, 2)
			So(adminUser.Groups().Ids(), ShouldContain, adminGrp.ID())
			So(adminUser.Groups().Ids(), ShouldContain, everyoneGroup.ID())
		})
	})
}

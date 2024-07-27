package e2e_tests

//
//func TestGet(t *testing.T) {
//	conn, err := grpc.NewClient(
//		fmt.Sprintf("0.0.0.0:%d", cfg.GrpcPort),
//		grpc.WithTransportCredentials(insecure.NewCredentials()),
//	)
//	require.NoError(t, err)
//	defer conn.Close()
//
//	client := v1.NewTeamClient(conn)
//
//	testTeamID := "some-team-id"
//
//	resp, err := client.Get(context.Background(), &v1.GetTeamRequest{
//		TeamId: testTeamID,
//	})
//	require.NoError(t, err)
//	require.NotNil(t, resp)
//
//	t.Log(resp)
//}
//
//func TestGetUserTeams(t *testing.T) {
//	conn, err := grpc.NewClient(
//		fmt.Sprintf("0.0.0.0:%d", cfg.GrpcPort),
//		grpc.WithTransportCredentials(insecure.NewCredentials()),
//	)
//	require.NoError(t, err)
//	defer conn.Close()
//
//	client := v1.NewTeamClient(conn)
//
//	resp, err := client.GetUserTeams(context.Background(), &emptypb.Empty{})
//	require.NoError(t, err)
//	require.NotNil(t, resp)
//
//	t.Log(resp)
//}
//
//func TestGetAllCanJoin(t *testing.T) {
//	conn, err := grpc.NewClient(
//		fmt.Sprintf("0.0.0.0:%d", cfg.GrpcPort),
//		grpc.WithTransportCredentials(insecure.NewCredentials()),
//	)
//	require.NoError(t, err)
//	defer conn.Close()
//
//	client := v1.NewTeamClient(conn)
//
//	resp, err := client.GetAllCanJoin(context.Background(), &emptypb.Empty{})
//	require.NoError(t, err)
//	require.NotNil(t, resp)
//
//	t.Log(resp)
//}
//
//func TestCreate(t *testing.T) {
//	conn, err := grpc.NewClient(
//		fmt.Sprintf("0.0.0.0:%d", cfg.GrpcPort),
//		grpc.WithTransportCredentials(insecure.NewCredentials()),
//	)
//	require.NoError(t, err)
//	defer conn.Close()
//
//	client := v1.NewTeamClient(conn)
//
//	resp, err := client.Create(context.Background(), &v1.CreateTeamRequest{
//		Name:        "Test Team",
//		Description: "This is a e2e_tests team",
//	})
//	require.NoError(t, err)
//	require.NotNil(t, resp)
//
//	t.Log(resp)
//}
//
//func TestJoin(t *testing.T) {
//	conn, err := grpc.NewClient(
//		fmt.Sprintf("0.0.0.0:%d", cfg.GrpcPort),
//		grpc.WithTransportCredentials(insecure.NewCredentials()),
//	)
//	require.NoError(t, err)
//	defer conn.Close()
//
//	client := v1.NewTeamClient(conn)
//
//	resp, err := client.Join(context.Background(), &v1.JoinTeamRequest{
//		TeamId: "some-team-id",
//	})
//	require.NoError(t, err)
//	require.NotNil(t, resp)
//
//	t.Log(resp)
//}
//
//func TestGetRoles(t *testing.T) {
//	conn, err := grpc.NewClient(
//		fmt.Sprintf("0.0.0.0:%d", cfg.GrpcPort),
//		grpc.WithTransportCredentials(insecure.NewCredentials()),
//	)
//	require.NoError(t, err)
//	defer conn.Close()
//
//	client := v1.NewTeamClient(conn)
//
//	resp, err := client.GetRoles(context.Background(), &v1.GetTeamRolesRequest{
//		TeamId: "some-team-id",
//	})
//	require.NoError(t, err)
//	require.NotNil(t, resp)
//
//	t.Log(resp)
//}
//
//func TestChangeRole(t *testing.T) {
//	conn, err := grpc.NewClient(
//		fmt.Sprintf("0.0.0.0:%d", cfg.GrpcPort),
//		grpc.WithTransportCredentials(insecure.NewCredentials()),
//	)
//	require.NoError(t, err)
//	defer conn.Close()
//
//	client := v1.NewTeamClient(conn)
//
//	resp, err := client.ChangeRole(context.Background(), &v1.ChangeTeamRole{
//		RoleId: "some-role-id",
//		UserId: "some-user-id",
//	})
//	require.NoError(t, err)
//	require.NotNil(t, resp)
//
//	t.Log(resp)
//}
//
//func TestLeave(t *testing.T) {
//	conn, err := grpc.NewClient(
//		fmt.Sprintf("0.0.0.0:%d", cfg.GrpcPort),
//		grpc.WithTransportCredentials(insecure.NewCredentials()),
//	)
//	require.NoError(t, err)
//	defer conn.Close()
//
//	client := v1.NewTeamClient(conn)
//
//	resp, err := client.Leave(context.Background(), &v1.LeaveTeamRequest{
//		TeamId: "some-team-id",
//	})
//	require.NoError(t, err)
//	require.NotNil(t, resp)
//
//	t.Log(resp)
//}

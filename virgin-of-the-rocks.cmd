ffmpeg -y \
	-r 30 -t 36.80 -loop 1 -i "virgin-of-the-rocks-paris.jpg" \
	-r 1 -t 3.50 -loop 1 -i "/home/mb/gits/mkshort/.cache/90ed2c0bdd3b03acdd5ac3a66e78983c197933eac4d771a64cdd353d7f12542d.png" \
	-r 1 -t 3.00 -loop 1 -i "/home/mb/gits/mkshort/.cache/0874052b84eaea473405843666c658223881605489935376624fed4601c3ab50.png" \
	-r 1 -t 3.00 -loop 1 -i "/home/mb/gits/mkshort/.cache/24292a73f7cf7157b6b62f4799931e9c8fca518208eeda4782c8bb75ef36a2dc.png" \
	-r 1 -t 3.00 -loop 1 -i "/home/mb/gits/mkshort/.cache/2b96a193886a1a0a329a24bcd3468f98f13cb337b26b4045e5843b5563252b47.png" \
	-r 1 -t 3.00 -loop 1 -i "/home/mb/gits/mkshort/.cache/23ac3136c67210652b5582c7847a47787f68f50e2288d4381b92d1607489e42d.png" \
	-r 1 -t 5.00 -loop 1 -i "/home/mb/gits/mkshort/.cache/9ec644c436c5e7a339ed857abc939e27b3f42362f21a62288e7c616ef4f3ce63.png" \
	-r 1 -t 3.00 -loop 1 -i "/home/mb/gits/mkshort/.cache/242215143f65b91164cfe12d9ba126b7fe52bf35eab7df8de227f1916e10e8b7.png" \
	-r 1 -t 4.00 -loop 1 -i "/home/mb/gits/mkshort/.cache/c9ed7cf3d1331df8d11fccf2d964fe221b1e037980e2346dd579604b94d9c53e.png" \
	-r 30 -t 17.30 -loop 1 -i "virgin-of-the-rocks-london.jpg" \
	-r 1 -t 4.00 -loop 1 -i "/home/mb/gits/mkshort/.cache/3d426114dc910999e9f5c9583b24444111dffda5ad5acf4e711d7079c83b4008.png" \
	-r 1 -t 2.50 -loop 1 -i "/home/mb/gits/mkshort/.cache/87f8eab5e96fcfde670fb93dfaa0649ba964286be63314c8fec524d1eafc184b.png" \
	-r 1 -t 3.00 -loop 1 -i "/home/mb/gits/mkshort/.cache/7f55fb922a6b776d224d061a9fcf1c725833680970f7d652965c735fbdde4342.png" \
	-filter_complex "
		[0:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in0];
		[9:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in1];

		[in0] [1:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.70,4.20)' [in0,1];
		[in0,1] [2:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,4.90,7.90)' [in0,2];
		[in0,2] [3:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,8.60,11.60)' [in0,3];
		[in0,3] [4:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,12.30,15.30)' [in0,4];
		[in0,4] [5:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,16.00,19.00)' [in0,5];
		[in0,5] [6:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,19.70,24.70)' [in0,6];
		[in0,6] [7:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,25.40,28.40)' [in0,7];
		[in0,7] [8:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,31.80,35.80)' [in0,9];
		[in1] [10:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.70,4.70)' [in1,1];
		[in1,1] [11:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,8.10,10.60)' [in1,3];
		[in1,3] [12:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,11.30,14.30)' [in1,4];

		[in0,9] [in1,4] concat=n=2:v=1:a=0:unsafe=1 [v]
	" \
	-pix_fmt yuv420p -r 30 -movflags faststart -map "[v]" "/dev/stdout"

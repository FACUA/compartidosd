package app

import "facua.org/compartidosd/client/fs"

func mountShares(shares []share) (mounted []share, errored []share) {
	for _, share := range shares {
		err := fs.MountShare(share)

		if err == nil {
			mounted = append(mounted, share)
		} else {
			errored = append(errored, share)
		}
	}

	return mounted, errored
}
